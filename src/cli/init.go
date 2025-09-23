/*
   Truct, Pretty minimal workflow runner.
   Open-Source, WTFPL License.

   Copyright (C) 2025-20xx Neo <ntsk@airmail.cc>
*/

package Cli

import (
    "fmt"
    "os"
	Types "github.com/neotesk/truct/src/types"
	Internal "github.com/neotesk/truct/src/util"
);

var Init = Types.Command {
    Name: "init",
    ShortDesc: "Creates a workflow file",
    LongDesc: "This command will create a workflow file, the file will contain a shell action and some project details. Don't forget to edit the newly created workflow file.",
    Descriptor: "",
    Action: func ( args []string, flags map[ string ] bool, adjectives map[ string ] string ) {
        tructFileName := adjectives[ "filename" ];
        doesTructFileExists := Internal.HandleError( Internal.FileSystem.Exists( tructFileName ) );
        if doesTructFileExists {
            Internal.ErrPrintf( "Workflow file '%s' already exists!\n", tructFileName );
            os.Exit( 1 );
        }

        name := "";
        fmt.Printf( "%s: ", Internal.Boldify( "Name" ) )
        fmt.Scanln( &name );

        description := "";
        fmt.Printf( "%s: ", Internal.Boldify( "Description" ) )
        fmt.Scanln( &description );

        version := "";
        fmt.Printf( "%s: ", Internal.Boldify( "Version" ) )
        fmt.Scanln( &version );

        output := []byte( "# Visit https://github.com/neotesk/truct/wiki for\n# more information.\n\nproject:\n    name: " + name + "\n    description: " + description + "\n    version: " + version + "\n\n.setup:\n    description: Performs a setup operation\n    actions:\n        -\n            action: shell\n            cmdlines:\n                - echo setup is done.\n\n.build:\n    description: Performs a build operation\n    actions:\n        -\n            action: shell\n            cmdlines:\n                - echo Build is done." );
        Internal.HandleErrorVoid( os.WriteFile( tructFileName, output, os.ModePerm ) );
    },
};