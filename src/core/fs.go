/*
   Truct, Pretty minimal workflow manager.
   Open-Source, WTFPL License.

   Copyright (C) 2025-20xx Neo <neotesk@airmail.cc>
*/

package Core

import (
	"fmt"
	"os"
	"path"

	Types "github.com/neotesk/truct/src/types"
	Internal "github.com/neotesk/truct/src/util"
)

var Move = Types.Action {
    Name: "move",
    Expects: map[ string ] any {
        "source": "src",
        "destination": "build",
        "overwrite": true,
    },
    Action: func( cwd string, details map[ string ] any, params Types.TructWorkflowRunArgs ) error {
        src := Internal.Make[ string ]( details[ "source" ] );
        dst := Internal.Make[ string ]( details[ "destination" ] );
        if ( !params.CommandLineArgs.TructFile.Settings.Silent ) {
            fmt.Printf( "%s Moving %s to %s\n", Internal.Colorify( "|", "3e83d6" ), Internal.Colorify( src, "ada440" ), Internal.Colorify( dst, "ada440" ) );
        }
        return Internal.FileSystem.Move(
            path.Join( cwd, src ),
            path.Join( cwd, dst ),
            Internal.Make[ bool ]( details[ "overwrite" ] ),
        )
    },
}

var Copy = Types.Action {
    Name: "copy",
    Expects: map[ string ] any {
        "source": "src",
        "destination": "build",
        "overwrite": true,
        "preservePermissions": true,
    },
    Action: func( cwd string, details map[ string ] any, params Types.TructWorkflowRunArgs ) error {
        src := Internal.Make[ string ]( details[ "source" ] );
        dst := Internal.Make[ string ]( details[ "destination" ] );
        shouldThrow := !Internal.Make[ bool ]( details[ "skipOnError" ] );
        if ( !params.CommandLineArgs.TructFile.Settings.Silent ) {
            fmt.Printf( "%s Copying %s to %s\n", Internal.Colorify( "|", "3e83d6" ), Internal.Colorify( src, "ada440" ), Internal.Colorify( dst, "ada440" ) );
        }

        output := Internal.FileSystem.Copy(
            path.Join( cwd, src ),
            path.Join( cwd, dst ),
            Internal.Make[ bool ]( details[ "overwrite" ] ),
            Internal.Make[ bool ]( details[ "preservePermissions" ] ),
        );

        if shouldThrow {
            return output;
        } else {
            return nil;
        }
    },
}

var Mkdir = Types.Action {
    Name: "mkdir",
    Expects: map[ string ] any {
        "destination": "build",
    },
    Action: func( cwd string, details map[ string ] any, params Types.TructWorkflowRunArgs ) error {
        dst := Internal.Make[ string ]( details[ "destination" ] );
        shouldThrow := !Internal.Make[ bool ]( details[ "skipOnError" ] );
        if ( !params.CommandLineArgs.TructFile.Settings.Silent ) {
            fmt.Printf( "%s Creating directory at %s\n", Internal.Colorify( "|", "3e83d6" ), Internal.Colorify( dst, "ada440" ) );
        }
        output := os.MkdirAll(
            path.Join( cwd, dst ),
            os.ModePerm.Perm(),
        );

        if shouldThrow {
            return output;
        } else {
            return nil;
        }
    },
}

var Touch = Types.Action {
    Name: "touch",
    Expects: map[ string ] any {
        "destination": "build",
        "content": "",
    },
    Action: func( cwd string, details map[ string ] any, params Types.TructWorkflowRunArgs ) error {
        dst := Internal.Make[ string ]( details[ "destination" ] );
        shouldThrow := !Internal.Make[ bool ]( details[ "skipOnError" ] );
        if ( !params.CommandLineArgs.TructFile.Settings.Silent ) {
            fmt.Printf( "%s Creating file at %s\n", Internal.Colorify( "|", "3e83d6" ), Internal.Colorify( dst, "ada440" ) );
        }
        output := os.WriteFile(
            path.Join( cwd, dst ),
            []byte( Internal.Make[ string ]( details[ "content" ] ) ),
            os.ModePerm.Perm(),
        );

        if shouldThrow {
            return output;
        } else {
            return nil;
        }
    },
}