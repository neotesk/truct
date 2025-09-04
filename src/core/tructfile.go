/*
   Truct, Pretty minimal workflow manager.
   Open-Source, WTFPL License.

   Copyright (C) 2025-20xx Neo <neotesk@airmail.cc>
*/

package Core

import (
	"maps"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"slices"
	"strings"
	"time"

	"github.com/goccy/go-yaml"
	Types "github.com/neotesk/truct/src/types"
	Internal "github.com/neotesk/truct/src/util"
);


func ReadYAML ( filePath string ) ( map[ string ] any, error ) {
    absPath, _err := filepath.Abs( filePath );
    if _err != nil {
        return nil, _err;
    }
    bytes, _err := os.ReadFile( absPath );
	if _err != nil {
	    return nil, _err;
	}
	output := map[ string ] any{};
	_err = yaml.Unmarshal( bytes, &output );
	if _err != nil {
	    return nil, _err;
	}
	return output, nil;
}

func ReadTructFile ( filePath string, silent bool ) Types.TructFile {
    data := Internal.HandleError( ReadYAML( filePath ) );

    // Let's format the data to something that
    // we can use.
    validKeys := []string{ "project", "settings", "env", "var" };

    tructFile := Types.TructFileRaw {
        Project: Types.ProjectDetails {},
        ProjectRaw: map[ string ]any {},
        Settings: map[ string ]any {},
        Environment: map[ string ] any {},
        Variables: map[ string ] any {},
        Workflows: map[ string ] Types.Workflow {},
    }

    for key := range data {
        if ( strings.HasPrefix( key, "." ) ) {
            // This is a workflow definition
            dataAsMap := map[ string ] any{};
            wfName := strings.TrimPrefix( key, "." );
            ok := false;
            if dataAsMap, ok = data[ key ].( map[ string ] any ); !ok {
                Internal.ErrPrintf( "Fatal Error: Value of '%s' is not a valid workflow definition.\n", key );
                os.Exit( 5 );
            }
            Internal.MakeSure( dataAsMap[ "actions" ], "Workflow with the name '" + wfName + "' doesn't have an actions field." );
            expectsField := Internal.MakeCoalesce( dataAsMap[ "expects" ], []any {} );
            formattedExpects := []Types.WorkflowVariable {};
            amiOptional := false;
            for _, x := range expectsField {
                errMsg := "Workflow with the name '" + wfName + "' contains a malformed expects list. (";
                v := Internal.Make[ map[ string ] any ]( x );
                Internal.MakeSure( v[ "name" ], errMsg + "Doesn't contain 'name' key)" );
                Internal.MakeSure( v[ "optional" ], errMsg + "Doesn't contain 'optional' key)" );

                // Perform a tiny validation
                name := Internal.Make[ string ]( v[ "name" ] );
                optional := Internal.Make[ bool ]( v[ "optional" ] );
                if ( optional == true ) {
                    Internal.MakeSure( v[ "value" ], errMsg + "Doesn't contain 'value' key)" );
                    amiOptional = true;
                } else if _, ok := v[ "value" ]; ok {
                    Internal.ErrPrintf( "Fatal Error: Expect Variable with the name '%s' cannot have a default value (non-optional expect) in workflow '%s'\n", name, wfName );
                    os.Exit( 1 );
                } else if amiOptional && !optional {
                    Internal.ErrPrintf( "Fatal Error: Expect Variable with the name '%s' cannot be a required value (required expect after optional expect) in workflow '%s'\n", name, wfName );
                    os.Exit( 1 );
                }

                output := Types.WorkflowVariable {
                    Name: name,
                    Value: Internal.MakeCoalesce( v[ "value" ], "" ),
                    Optional: optional,
                }
                formattedExpects = append( formattedExpects, output )
            }

            result := Types.Workflow {
                Description: Internal.MakeCoalesce( dataAsMap[ "description" ], "No description provided." ),
                Actions: Internal.Make[ []any ]( dataAsMap[ "actions" ] ),
                Expects: formattedExpects,
            }

            tructFile.Workflows[ wfName ] = result;
        } else if slices.Contains( validKeys, key ) {
            // This is a project key.
            obj := Internal.Make[ map[ string ] any ]( data[ key ] );
            switch key {
                case "project":
                    tructFile.Project.Name = Internal.MakeCoalesce( obj[ "name" ], "Unknown" );
                    tructFile.Project.Description = Internal.MakeCoalesce( obj[ "description" ], "Unknown" );
                    tructFile.Project.Version = Internal.MakeCoalesce( obj[ "version" ], "Unknown" );
                    tructFile.Project.Repository = Internal.MakeCoalesce( obj[ "repository" ], "Unknown" );
                    tructFile.ProjectRaw = obj;
                case "settings":
                    tructFile.Settings = obj;
                case "env":
                    tructFile.Environment = obj;
                case "var":
                    tructFile.Variables = obj;
            }
        } else {
            // Seems like it's an invalid key
            Internal.ErrPrintf( "Fatal Error: Invalid key '%s' in workflow file.\n", key );
            os.Exit( 4 );
        }
    }

    output := Types.TructFile {
        Project: tructFile.Project,
        ProjectRaw: tructFile.ProjectRaw,
        Environment: tructFile.Environment,
        Variables: tructFile.Variables,
        Workflows: tructFile.Workflows,
        Settings: Types.TructSettings {},
    }

    // After we fetch the project details, we will now
    // create defaults and properly object-ify them
    defaultShell := "unknown"
    switch runtime.GOOS {
        case "windows":
            defaultShell = "cmd.exe /k"
        case "darwin", "linux", "android", "freebsd", "openbsd", "solaris", "plan9", "netbsd", "dragonfly":
            defaultShell = "sh -c"
        default:
            Internal.ErrPrintf( "Fatal Error: Unsupported Operating System (Default system shell cannot be determined)\n" );
            os.Exit( 4 );
    }

    settings := Internal.MakeCoalesce( tructFile.Settings[ "truct" ], map[ string ] any {
        "reportActions": nil,
        "silent": nil,
    } );

    osSettings := Internal.MakeCoalesce( tructFile.Settings[ "os." + runtime.GOOS ], map[ string ] any {
        "shell": nil,
    } );

    output.Settings.ReportActions = Internal.MakeCoalesce( settings[ "reportActions" ], true );
    output.Settings.Silent = Internal.MakeCoalesce( settings[ "silent" ], silent );
    output.Settings.Shell = Internal.MakeCoalesce( osSettings[ "shell" ], defaultShell );

    return output;
}

func CreateRootVarTable ( tructFile Types.TructFile ) map[ string ] string {
    varTable := map[ string ] string {};

    // OS Variables
    varTable[ "os.type" ] = runtime.GOOS;
    varTable[ "os.arch" ] = runtime.GOARCH;

    // Time
    varTable[ "rt.time" ] = time.Now().Format( time.RFC3339Nano );

    // Environment
    for _, env := range os.Environ() {
        fIndex := strings.Index( env, "=" );
        key := env[ :fIndex ];
        value := env[ fIndex + 1: ];
        varTable[ "env." + key ] = value;
    }

    // Project Details
    for key, value := range tructFile.ProjectRaw {
        varTable[ "project." + key ] = FormatVariables( Internal.Make[ string ]( value ), varTable );
    }

    // Environment Overrides
    for key, value := range tructFile.Environment {
        varTable[ "env." + key ] = FormatVariables( Internal.Make[ string ]( value ), varTable );
    }

    // Variables
    for key, value := range tructFile.Variables {
        varTable[ "var." + key ] = FormatVariables( Internal.Make[ string ]( value ), varTable );
    }

    return varTable;
}

var varFormatRegex = regexp.MustCompile( `{{([\w\s\d.$]+)}}` );
func FormatVariables ( input string, varTable map[ string ] string ) string {
    return varFormatRegex.ReplaceAllStringFunc( input, func ( stri string ) string {
        val, ok := varTable[ stri[ 2:len( stri ) - 2 ] ];
        if !ok { return stri; }
        return val;
    } );
}

func MType [ T any ] ( input any, errText string ) T {
    if m, o := input.( T ); o {
        return m;
    }
    Internal.ErrPrintf( "Fatal Error: %s\n", errText );
    os.Exit( 4 );
    return input.( T );
}

func WorkflowScope ( wf Types.Workflow, varTable map[ string ] string, argList []string ) map[ string ] string {
    scope := map[ string ] string {};
    maps.Insert( scope, maps.All( varTable ) );
    if len( wf.Expects ) != 0 {
        // Yep, there's a descriptor
        containsOptional := false;
        for i, val := range wf.Expects {
            composite := Types.WorkflowVariable {
                Name: MType[ string ]( val.Name, "'name' must be a string" ),
                Optional: MType[ bool ]( val.Optional, "'optional' must be a boolean" ),
            };
            if containsOptional && !composite.Optional {
                Internal.ErrPrintf( "Fatal Error: Expects list of workflow '%s' has a required expection after optional expection.\n", wf.Name );
                os.Exit( 4 );
            };
            if composite.Optional {
                containsOptional = true;
                composite.Value = FormatVariables(
                    MType[ string ]( val.Value, "'value' must be a string" ),
                    scope,
                );
            } else {
                if i >= len( argList ) {
                    Internal.ErrPrintf( "Fatal Error: Missing argument for '%s' as index %d\n", composite.Name, i );
                    os.Exit( 4 );
                }
            }
            if i < len( argList ) { composite.Value = argList[ i ]; }
            scope[ "$." + composite.Name ] = composite.Value;
        }
    }
    return scope;
}

func MarshalVar ( thing any, scope map[ string ] string ) any {
    if det, ok := thing.( map[ string ] any ); ok {
        output := map[ string ] any {};
        for key := range det {
            output[ key ] = MarshalVar( det[ key ], scope );
        }
        return output;
    } else if det, ok := thing.( string ); ok {
        return FormatVariables( det, scope );
    } else if det, ok := thing.( []any ); ok {
        outp := make( []any, len( det ) );
        for i := range det {
            outp[ i ] = MarshalVar( det[ i ], scope );
        }
        return outp;
    }
    return thing;
}

func RunAction ( actionDetails map [ string ] any, custDetails Types.TructWorkflowRunArgs ) {
    actName := Internal.Make[ string ]( actionDetails[ "action" ] );
    act := GetAction( actName );

    aData := Internal.CopyMap( act.Expects );
    maps.Copy( aData, actionDetails );

    if _, ok := aData[ "skipOnError" ]; !ok {
        aData[ "skipOnError" ] = false;
    }

    Internal.HandleErrorVoid( act.Action( custDetails.WorkingDirectory, aData, custDetails ) );
}

func RunWorkflow ( wfMain Types.TructWorkflowRunArgs ) {
    curWf, ok := wfMain.CommandLineArgs.TructFile.Workflows[ wfMain.WorkflowName ];
    if !ok {
        Internal.ErrPrintf( "Fatal Error: Workflow with the name '%s' does not exist.\n", wfMain.WorkflowName );
        os.Exit( 4 );
    }

    curScope := WorkflowScope( curWf, wfMain.ScopeVariables, wfMain.CommandLineArgs.Keywords );

    for _, action := range curWf.Actions {
        if actionCmd, ok := action.( string ); ok {
            // It's a cmd/shortcut/alias (whatever
            // you wanna call it)
            cmd := Internal.ParseCmdline( FormatVariables( actionCmd, curScope ) );
            if len( cmd ) < 1 {
                Internal.ErrPrintf( "Fatal Error: Workflow with the name '%s' has an incomplete alias call inside actions list.\n", wfMain.WorkflowName );
                os.Exit( 4 );
            }
            cust := Types.TructWorkflowRunArgs {
                WorkflowName: cmd[ 0 ],
                ScopeVariables: curScope,
                CommandLineArgs: Types.CommandLineArgs {
                    Keywords: cmd[ 1: ],
                    Arguments: wfMain.CommandLineArgs.Arguments,
                    Flags: wfMain.CommandLineArgs.Flags,
                    TructFile: wfMain.CommandLineArgs.TructFile,
                },
                WorkingDirectory: wfMain.WorkingDirectory,
            }
            RunWorkflow( cust );
            continue;
        }

        scoped := Internal.Make[ map[ string ] any ]( MarshalVar( action, curScope ) );
        RunAction( scoped, wfMain );
    }
}