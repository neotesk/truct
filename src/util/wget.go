/*
   Truct, Pretty minimal workflow runner.
   Open-Source, WTFPL License.

   Copyright (C) 2025-20xx Neo <ntsk@airmail.cc>
*/

package Internal

import (
    "fmt"
    "io"
    "net/http"
    "os"
    "time"
)

type ProgressReader struct {
    Reader io.Reader;
    Size   int64;
    Pos    int64;
    Report bool;
    URL    string;
}

func ( pr *ProgressReader ) Read ( p []byte ) ( int, error ) {
    n, err := pr.Reader.Read( p );
    if err == nil {
        pr.Pos += int64( n );
        if !pr.Report { return n, err; }
        if pr.Size == -1 {
            fmt.Printf( "\r%s Downloading to %s " + Colorify( "%d Bytes ", "9f9f9f" ), Colorify( "-", "3e83d6" ), Colorify( pr.URL, "ada440" ), pr.Pos );
        } else {
            fmt.Printf( "\r%s Downloading to %s " + Colorify( "%.2f%% ", "9f9f9f" ), Colorify( "-", "3e83d6" ), Colorify( pr.URL, "ada440" ), float64( pr.Pos ) / float64( pr.Size ) * 100 );
        }

    }
    return n, err;
}

func Download ( url string, path string, reportStatus bool, reportUri string, shouldReportError bool ) error {
    start := time.Now().UnixMilli();
    req, _ := http.NewRequest( "GET", url, nil );
    resp, _ := http.DefaultClient.Do( req );
    if resp.StatusCode != 200 && shouldReportError {
        ErrPrintf( "| Error while downloading: %v\n", resp.StatusCode );
        return fmt.Errorf( "Error while downloading: %v", resp.StatusCode );
    }
    defer resp.Body.Close();

    f, _ := os.OpenFile( path, os.O_CREATE | os.O_WRONLY, 0644 );
    defer f.Close();

    progressReader := &ProgressReader{
        Reader: resp.Body,
        Size:   resp.ContentLength,
        Report: reportStatus,
        URL:    reportUri,
    }

    if _, err := io.Copy( f, progressReader ); err != nil && shouldReportError {
        ErrPrintf( "| Error while downloading: %v\n", err );
        return err;
    }

    if reportStatus {
        fmt.Printf( "\r\x1b[K%s Downloaded to %s, Took %.2fs\n", Colorify( "|", "3e83d6" ), reportUri, float64( time.Now().UnixMilli() - start ) / 1000 );
    }

    return nil;
}