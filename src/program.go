/*
   Truct, Pretty minimal workflow manager.
   Open-Source, WTFPL License.

   Copyright (C) 2025-20xx Neo <neotesk@airmail.cc>
*/

package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	Cli "github.com/neotesk/truct/src/cli"
	Core "github.com/neotesk/truct/src/core"
	Types "github.com/neotesk/truct/src/types"
	Internal "github.com/neotesk/truct/src/util"
)

// This will be changed during build
var versionString = "github-release";

func StringChunk ( str string, totalLength int ) string {
    chunks := strings.Split( str, " " );
    newText := "";
    currentText := []string{};
    for _, chunk := range chunks {
        currentText = append( currentText, chunk );
        if len( strings.Join( currentText, " " ) ) > totalLength {
            newText = newText + strings.Join( currentText[ :len( currentText ) - 1 ], " " ) + "\n";
            currentText = []string{ chunk };
        }
    }
    return newText + strings.Join( currentText, " " );
}

func FindLongestCmdLength ( cmds []Types.Command ) int {
    currLength := 0;
    for _, text := range cmds {
        tLen := len( text.Name );
        if ( tLen > currLength ) {
            currLength = tLen;
        }
    }
    return currLength + 3;
}

func FindFlag ( flagTag string, flagList []Types.Flag ) ( *Types.Flag, bool ) {
    for _, Flag := range flagList {
        if Flag.Name == flagTag {
            return &Flag, true;
        }
    }
    return nil, false;
}

func RunProgram ( prog Types.Program, args Types.ArgumentsList ) {
    // First, let's get the command details
    command := "";
    commandArgs := []string{};
    if len( args.Keywords ) != 0 {
        command = args.Keywords[ 0 ];
        commandArgs = args.Keywords[ 1: ];
    }

    // After that we will compile the program
    // into a new object
    commandsList := map[ string ] Types.Command {};
    for _, command := range prog.Commands {
        commandsList[ command.Name ] = command;
    }

    // After compiling the commands to a list
    // we will perform a search
    commandDetails, commandExists := commandsList[ command ];
    if commandExists {
        commandDetails.Action( commandArgs, args.Flags, args.Arguments );
        return;
    }

    // Display version if possible
    if command == "version" || args.Flags[ "v" ] {
        fmt.Println( versionString );
        return;
    }

    // Seems like the command doesn't exist, check
    // if there's any builtin commands.
    if command == "help" || command == "" {
        // Create a list of flags and arguments
        flags := "";
        arguments := []string{};
        for _, key := range prog.DefaultArgs.Flags {
            flags += key.Name;
        }
        for key := range args.Arguments {
            arguments = append( arguments, "--" + key );
        }
        fmt.Println( Internal.Colorify( StringChunk( prog.Desc, 60 ), "ab86c2") );

        // Before showing the help page, we need to check if
        // there's any keyword after our command, and so we
        // must check.
        argsLen := len( args.Keywords );
        if argsLen > 1 && args.Keywords[ 0 ] == "help" {
            firstArg := args.Keywords[ 1 ];
            if argsLen == 3 && firstArg == "do" {
                // This is a workflow help page.
                tructFileName := args.Arguments[ "filename" ];
                wfName := args.Keywords[ 2 ];

                // Check if the truct file exists.
                doesTructFileExists := Internal.HandleError( Internal.FileSystem.Exists( tructFileName ) );
                if !doesTructFileExists {
                    Internal.ErrPrintf( "Workflow file '%s' doesn't exist, dummy!\n", tructFileName );
                    os.Exit( 1 );
                }

                // Read the file.
                tructFile := Core.ReadTructFile( tructFileName, true );
                wf := tructFile.Workflows[ wfName ];

                xargs := []string {}
                for _, wv := range wf.Expects {
                    outp := "<%s>";
                    if wv.Optional {
                        outp = "[%s]";
                    }
                    xargs = append( xargs, fmt.Sprintf( outp, wv.Name ) );
                }

                fmt.Printf( "\n%s: %s %s %s\n\n%s\n",
                    Internal.Colorify( "Usage", "fcba03"),
                    Internal.Colorify( "truct", "2190bf"),
                    Internal.Colorify( "do " + wfName, "57879c"),
                    Internal.Colorify( strings.Join( xargs, " " ), "7a7a7a"),
                    StringChunk( wf.Description, 60 ) );
                os.Exit( 0 );
            } else if argsLen == 2 && firstArg == "flags" {
                // This is a flags help page
                fmt.Printf( "\n%s\n", Internal.Colorify( "Flags:", "32ab73") );
                for _, flag := range prog.DefaultArgs.Flags {
                    fmt.Printf( "   -%s   %s\n", flag.Name, flag.ShortDesc );
                }
                os.Exit( 0 );
            } else if argsLen == 2 {
                // This is a command help page
                Command := Cli.GetCommandByName( firstArg );
                if Command == nil {
                    Internal.ErrPrintf( "\nCommand with the name '%s' does not exist.\n", firstArg );
                    os.Exit( 1 );
                }
                fmt.Printf( "\n%s: %s %s %s\n\n%s\n",
                    Internal.Colorify( "Usage", "fcba03"),
                    Internal.Colorify( "truct", "2190bf"),
                    Internal.Colorify( Command.Name, "57879c"),
                    Internal.Colorify( Command.Descriptor, "7a7a7a"),
                    StringChunk( Command.LongDesc, 60 ) );
                os.Exit( 0 );
            }
        }

        fmt.Printf( "\n%s:\n    %s [%s] [%s] %s %s\n",
            Internal.Colorify( "Usage", "fcba03"),
            Internal.Colorify( prog.Name, "2190bf"),
            Internal.Colorify( "-" + flags, "7a7a7a"),
            Internal.Colorify( strings.Join( arguments, ", " ), "7a7a7a"),
            Internal.Colorify( "<command>", "9f9f9f"),
            Internal.Colorify( "[arguments]", "9f9f9f") );

        // Now we're gonna list commands, however we
        // need to find out the longest command so we
        // can create a padding.
        fmt.Printf( "\n%s:\n", Internal.Colorify( "Commands", "32ab73") );
        padLen := FindLongestCmdLength( prog.Commands );

        // We will put a pseudo-command called "help"
        fmt.Printf( "    %-" + strconv.Itoa( padLen ) + "s%s\n", "help", Internal.Colorify( "Shows the help page", "9f9f9f") );

        for _, command := range prog.Commands {
            fmt.Printf( "    %-" + strconv.Itoa( padLen ) + "s%s\n", command.Name, Internal.Colorify( command.ShortDesc, "9f9f9f") );
        }

        // After that, we will show that there's actually more
        // help pages available
        fmt.Printf( "\n%s:\n", Internal.Colorify( "More help commands", "32ab73") );
        fmt.Println( "    help flags             " + Internal.Colorify( "Lists the arguments and flags", "9f9f9f") );
        fmt.Println( "    help [command]         " + Internal.Colorify( "Shows a help page about a command", "9f9f9f") );
        fmt.Println( "    help do [workflow]     " + Internal.Colorify( "Describes a workflow", "9f9f9f") );

        fmt.Println( "\n" + Internal.Colorify( StringChunk( prog.Footer, 60 ), "c96d7f") );
        return;
    }

    // Seems like there is nothing...
    Internal.ErrPrintf( "Unknown command: %s\nPlease visit the help page using the command down below.\n   truct help\n", command );
}