package app

import (

    "io"
    "mime/multipart"

)
// readFile reads the content of an uploaded file
func readFile(file *multipart.FileHeader) (string, error) {
    f, err := file.Open()
    if err != nil {
        return "", err
    }
    defer f.Close()

    content, err := io.ReadAll(f)
    if err != nil {
        return "", err
    }

    return string(content), nil
}