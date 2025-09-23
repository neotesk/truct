/*
   Truct, Pretty minimal workflow runner.
   Open-Source, WTFPL License.

   Copyright (C) 2025-20xx Neo <ntsk@airmail.cc>
*/

package Internal

import (
	"io"
	"os"
	"path"
	"path/filepath"
)

// Filesystem Namespace
type IFileSystem struct {};
var FileSystem IFileSystem;

func ( IFileSystem ) Exists ( path string ) ( bool, error ) {
    absPath, _err := filepath.Abs( path );
    if _err != nil {
        return false, _err;
    }
    _, _err = os.Stat( absPath );
    return _err == nil, nil;
}

func copyFile ( dst, src string, overwrite bool ) error {
    srcF, err := os.Open( src );
    if err != nil {
        return err;
    }
    defer srcF.Close();
    info, err := srcF.Stat();
    if err != nil {
        return err;
    }
    var dstF *os.File;
    if overwrite {
        dstF, err = os.OpenFile( dst, os.O_RDWR|os.O_CREATE|os.O_TRUNC, info.Mode() );
        if err != nil {
            return err;
        }
    } else {
        dstF, err = os.Create( dst );
        if err != nil {
            return err;
        }
        err = dstF.Chmod( info.Mode() );
        if err != nil {
            return err;
        }
    }
    defer dstF.Close();
    if _, err := io.Copy( dstF, srcF ); err != nil {
        return err;
    }
    return nil;
}

func ( IFileSystem ) Copy ( src, dst string, overwrite bool, preservePermissions bool ) error {
    srcStat, err := os.Stat( src );
    if err != nil {
        return err;
    }
    if !srcStat.IsDir() {
        return copyFile( dst, src, overwrite );
    }
    if overwrite {
        if err := os.RemoveAll( dst ); err != nil {
            return err;
        }
    }
    if err := os.CopyFS( dst, os.DirFS( src ) ); err != nil {
        return err;
    }
    if preservePermissions {
        if err := os.Chmod( dst, srcStat.Mode().Perm() ); err != nil {
            return err;
        }
    }
    return nil;
}

func ( IFileSystem ) Move ( src, dst string, overwrite bool ) error {
    if _, err := os.Stat( src ); err != nil {
        return err;
    }
    if stat, err := os.Stat( dst ); err == nil && stat.IsDir() && overwrite {
        fs, err := os.ReadDir( src );
        if err != nil {
            return err;
        }
        for _, dirEnt := range fs {
            newSrc := path.Join( src, dirEnt.Name() );
            newDst := path.Join( dst, dirEnt.Name() );
            err := FileSystem.Move( newSrc, newDst, overwrite );
            if err != nil {
                return err
            }
        }
        return os.Remove( src );
    } else if err := os.Rename( src, dst ); err != nil {
        return err;
    }
    return nil;
}