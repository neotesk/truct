/*
   Truct, Pretty minimal workflow manager.
   Open-Source, WTFPL License.

   Copyright (C) 2025-20xx Neo <neotesk@airmail.cc>
*/

package Core

import (
	"os"
	"path"

	Types "github.com/neotesk/truct/src/types"
	Internal "github.com/neotesk/truct/src/util"
)

var Get = Types.Action {
    Name: "get",
    Expects: map[ string ] any {
        "destination": "build.txt",
        "url": "",
        "skipIfExists": false,
        "extract": false,
    },
    Action: func( cwd string, details map[ string ] any, params Types.TructWorkflowRunArgs ) error {
        dst := Internal.Make[ string ]( details[ "destination" ] );
        ext := Internal.Make[ bool ]( details[ "extract" ] );
        shouldThrow := !Internal.Make[ bool ]( details[ "skipOnError" ] );
        if b, _ := Internal.FileSystem.Exists( dst ); Internal.Make[ bool ]( details[ "skipIfExists" ] ) && b {
            return nil;
        }
        if ext {
            newDst := dst + ".tmp";
            Internal.Download(
                Internal.Make[ string ]( details[ "url" ] ),
                path.Join( cwd, newDst ),
                !params.CommandLineArgs.TructFile.Settings.Silent,
                dst,
                shouldThrow,
            );
            if shouldThrow {
                Internal.HandleErrorVoid( Internal.Unzip( newDst, dst ) );
            } else {
                Internal.Unzip( newDst, dst );
            }
            os.Remove( newDst )
        } else {
            Internal.Download(
                Internal.Make[ string ]( details[ "url" ] ),
                path.Join( cwd, dst ),
                !params.CommandLineArgs.TructFile.Settings.Silent,
                dst,
                shouldThrow,
            );
        }
        return nil;
    },
}