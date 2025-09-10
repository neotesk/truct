/*
   Truct, Pretty minimal workflow manager.
   Open-Source, WTFPL License.

   Copyright (C) 2025-20xx Neo <ntsk@airmail.cc>
*/

package Core

import (
	"os"
	"path"

	Types "github.com/neotesk/truct/src/types"
	Internal "github.com/neotesk/truct/src/util"
)

var actions = []Types.Action {
    Shell,
    Copy,
    Move,
    Mkdir,
    Touch,
    Remove,
    Get,
    Zip,
    Unzip,
}

func actionsToMap () map [ string ] Types.Action {
    cib := map [ string ] Types.Action {};
    for _, a := range actions {
        cib[ a.Name ] = a;
    }
    return cib;
}

var actList = actionsToMap();

func GetAction ( name string ) Types.Action {
    act, ok := actList[ name ];

    // This is stupidly dirty but I had no other choice
    // since it warns me about initialization cycle.
    // If you have a better approach please fix this code.
    if name == "truct-run" {
        ok = true
        act = Types.Action {
            Name: "truct-run",
            Expects: map[ string ] any {
                "source": "",
                "workflow": "",
                "details": "",
            },
            Action: func( cwd string, details map[ string ] any, params Types.TructWorkflowRunArgs ) error {
                // Get the filename.
                tructFileName := Internal.Make[ string ]( details[ "source" ] );

                // Check if the truct file exists.
                doesTructFileExists := Internal.HandleError( Internal.FileSystem.Exists( tructFileName ) );
                if !doesTructFileExists {
                    Internal.ErrPrintf( "Workflow file '%s' doesn't exist, dummy!\n", tructFileName );
                    os.Exit( 1 );
                }

                // Read the file.
                tructFile := ReadTructFile( tructFileName, params.CommandLineArgs.TructFile.Settings.Silent );

                // Get the workflow name
                wfName := Internal.MakeCoalesce( details[ "workflow" ], "build" );

                // Now we will format variables based on environment
                // variables, OS details and much more.
                // Let's start with OS details first.
                varTable := CreateRootVarTable( tructFile, params.WorkingDirectory );

                RunWorkflow( Types.TructWorkflowRunArgs {
                    WorkflowName: wfName,
                    ScopeVariables: varTable,
                    CommandLineArgs: Types.CommandLineArgs{
                        Arguments: params.CommandLineArgs.Arguments,
                        Flags: params.CommandLineArgs.Flags,
                        Keywords: Internal.ParseCmdline( Internal.Make[ string ]( details[ "details" ] ) ),
                        TructFile: tructFile,
                    },
                    WorkingDirectory: path.Dir( tructFileName ),
                } );
                return nil;
            },
        }
    }

    if !ok {
        Internal.ErrPrintf( "Fatal Error: Action with the name '%s' does not exist.\n", name );
        os.Exit( 1 );
    }
    return act;
}