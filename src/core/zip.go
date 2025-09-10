/*
   Truct, Pretty minimal workflow manager.
   Open-Source, WTFPL License.

   Copyright (C) 2025-20xx Neo <ntsk@airmail.cc>
*/

package Core

import (
	"fmt"
	"os"

	Types "github.com/neotesk/truct/src/types"
	Internal "github.com/neotesk/truct/src/util"
)

var Unzip = Types.Action {
    Name: "unzip",
    Expects: map[ string ] any {
        "source": "src.zip",
        "destination": "build",
        "overwrite": true,
    },
    Action: func( cwd string, details map[ string ] any, params Types.TructWorkflowRunArgs ) error {
        src := Internal.Make[ string ]( details[ "source" ] );
        dst := Internal.Make[ string ]( details[ "destination" ] );
        pset := params.CommandLineArgs.TructFile.Settings;
        if ( !pset.Silent && pset.ReportActions ) {
            fmt.Printf( "%s Unzipping %s to %s\n", Internal.Colorify( "|", "3e83d6" ), Internal.Colorify( src, "ada440" ), Internal.Colorify( dst, "ada440" ) );
        }
        if ( Internal.Make[ bool ]( details[ "overwrite" ] ) ) {
            os.RemoveAll( dst );
        }
        return Internal.Unzip( src, dst );
    },
}

var Zip = Types.Action {
    Name: "zip",
    Expects: map[ string ] any {
        "source": "src",
        "destination": "build.zip",
        "overwrite": true,
        "preserveRoot": false,
    },
    Action: func( cwd string, details map[ string ] any, params Types.TructWorkflowRunArgs ) error {
        src := Internal.Make[ string ]( details[ "source" ] );
        dst := Internal.Make[ string ]( details[ "destination" ] );
        proot := Internal.Make[ bool ]( details[ "preserveRoot" ] );
        pset := params.CommandLineArgs.TructFile.Settings;
        if ( !pset.Silent && pset.ReportActions ) {
            fmt.Printf( "%s Zipping %s to %s\n", Internal.Colorify( "|", "3e83d6" ), Internal.Colorify( src, "ada440" ), Internal.Colorify( dst, "ada440" ) );
        }
        if ( Internal.Make[ bool ]( details[ "overwrite" ] ) ) {
            os.Remove( dst );
        }
        return Internal.Zip( src, dst, proot );
    },
}