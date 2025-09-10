/*
   Truct, Pretty minimal workflow manager.
   Open-Source, WTFPL License.

   Copyright (C) 2025-20xx Neo <ntsk@airmail.cc>
*/

package Cli

import (
	"os"

	Core "github.com/neotesk/truct/src/core"
	Types "github.com/neotesk/truct/src/types"
	Internal "github.com/neotesk/truct/src/util"
);

var Do = Types.Command {
    Name: "do",
    ShortDesc: "Runs a workflow",
    LongDesc: "This command runs a workflow, by default it will run the 'build' workflow, however you can supply one more argument to this command to set which workflow should be run.",
    Descriptor: "[workflow] [...args]",
    Action: func ( args []string, flags map[ string ] bool, adjectives map[ string ] string ) {
        // Get the filename.
        tructFileName := adjectives[ "filename" ];

        // Check if the truct file exists.
        doesTructFileExists := Internal.HandleError( Internal.FileSystem.Exists( tructFileName ) );
        if !doesTructFileExists {
            Internal.ErrPrintf( "Workflow file '%s' doesn't exist, dummy!\n", tructFileName );
            os.Exit( 1 );
        }

        // Read the file.
        tructFile := Core.ReadTructFile( tructFileName, flags[ "s" ] );

        // Get the workflow name
        wfPossible := Internal.PossibleItem( args, 0 );
        if wfPossible != nil {
            args = args [ 1: ];
        }
        wfName := Internal.MakeCoalesce( wfPossible, "build" );

        // Now we will format variables based on environment
        // variables, OS details and much more.
        // Let's start with OS details first.
        cwd := Internal.HandleError( os.Getwd() );
        varTable := Core.CreateRootVarTable( tructFile, cwd );

        Core.RunWorkflow( Types.TructWorkflowRunArgs {
            WorkflowName: wfName,
            ScopeVariables: varTable,
            CommandLineArgs: Types.CommandLineArgs{
                Arguments: adjectives,
                Flags: flags,
                Keywords: args,
                TructFile: tructFile,
            },
            WorkingDirectory: cwd,
        } );
    },
};