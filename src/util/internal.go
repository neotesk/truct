/*
   Truct, Pretty minimal workflow manager.
   Open-Source, WTFPL License.

   Copyright (C) 2025-20xx Neo <ntsk@airmail.cc>
*/

package Internal

import (
    "archive/zip"
	"os"
	"path"
	"strings"
	"path/filepath"
	"io"
)

func HandleError [ T any ] ( thing T, err error ) T {
    if err != nil {
        ErrPrintf( "Fatal Error: %s\n", err.Error() );
        os.Exit( 1 );
    }
    return thing;
}

func HandleErrorVoid ( err error ) {
    if err != nil {
        ErrPrintf( "Fatal Error: %s\n", err.Error() );
        os.Exit( 1 );
    }
}

func ParseCmdline ( input string ) []string {
    chartable := strings.Split( input, "" );
    current := "";
    args := []string {};
    inQuotes := false;
    quoteType := "";

    i := 0;
    for i < len( chartable ) {
        c := chartable[ i ];
        if ( c == "\\" ) {
            i++;
            current += chartable[ i ];
        } else if ( inQuotes ) {
            if ( c == quoteType ) {
                inQuotes = false;
                quoteType = "";
                args = append( args, current );
                current = "";
            }
            current += c;
        } else if c == "\"" {
            inQuotes = true;
            quoteType = "\"";
        } else if c == "'" {
            inQuotes = true;
            quoteType = "'";
        } else if c == " " {
            if ( current == "" ) {
                continue;
            }
            args = append( args, current );
            current = "";
        } else {
            current += c;
        }
        i++;
    }
    if ( current != "" ) {
        args = append( args, current );
        current = "";
    }
    return args;
}

func Unzip ( source, dest string ) error {
    read, err := zip.OpenReader( source );
    if err != nil { return err; }
    defer read.Close();
    for _, file := range read.File {
        if file.Mode().IsDir() { continue; }
        open, err := file.Open();
        if err != nil { return err; }
        name := path.Join( dest, file.Name );
        os.MkdirAll( path.Dir( name ), os.ModePerm.Perm() );
        create, err := os.Create( name );
        if err != nil { return err; }
        defer create.Close();
        create.ReadFrom( open );
    }
    return nil;
}

func Zip ( source, dest string, preserveRoot bool ) error {
    file := HandleError( os.Create( dest ) );
    defer file.Close();

    writer := zip.NewWriter( file );
    defer writer.Close();

    newSource := source;
    if preserveRoot {
        newSource = filepath.Clean( filepath.Join( source, "../" ) );
    }

    walker := func( walkSrc string, info os.FileInfo, err error ) error {
        if err != nil {
            return err;
        }
        if info.IsDir() {
            return nil;
        }
        file, err := os.Open( walkSrc );
        if err != nil {
            return err;
        }
        defer file.Close();
        f, err := writer.Create( HandleError( filepath.Rel( newSource, walkSrc ) ) );
        if err != nil {
            return err;
        }
        _, err = io.Copy( f, file );
        if err != nil {
            return err;
        }
        return nil;
    }

    return filepath.Walk( source, walker );
}