/*
   Truct, Pretty minimal workflow manager.
   Open-Source, WTFPL License.

   Copyright (C) 2025-20xx Neo <neotesk@airmail.cc>
*/

package Core

import (
	"fmt"
	"os"
	"os/exec"

	Types "github.com/neotesk/truct/src/types"
	Internal "github.com/neotesk/truct/src/util"
)

var Shell = Types.Action {
    Name: "shell",
    Expects: map[ string ] any {
        "env": map[ string ] any {},
        "cmdlines": []string {},
        "onlyErrors": false,
    },
    Action: func( cwd string, details map[ string ] any, params Types.TructWorkflowRunArgs ) error {
        for _, val := range Internal.Make[ []any ]( details[ "cmdlines" ] ) {
            executor := Internal.ParseCmdline( params.CommandLineArgs.TructFile.Settings.Shell );
            line := Internal.Make[ string ]( val );
            cmdline := append( executor, line );
            if ( !params.CommandLineArgs.TructFile.Settings.Silent ) {
                fmt.Printf( "%s Invoking %s\n", Internal.Colorify( "|", "3e83d6" ), Internal.Colorify( line, "ada440" ) );
            }
            cmd := exec.Command( cmdline[ 0 ], cmdline[ 1: ]... );
            cmd.Dir = cwd;
            cmd.Env = os.Environ();
            for key, val2 := range details[ "env" ].( map[ string ] any ) {
                cmd.Env = append( cmd.Env, key + "=" + val2.( string ) )
            }
            if !params.CommandLineArgs.TructFile.Settings.Silent {
                cmd.Stdin = os.Stdin
                cmd.Stdout = os.Stdout
                cmd.Stderr = os.Stderr
            }
            err := cmd.Run()
            if err != nil && !Internal.Make[ bool ]( details[ "skipOnError" ] ) {
                return err;
            }
        }
        return nil;
    },
}
