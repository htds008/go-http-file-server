package serverHandler

import (
	"../util"
	"errors"
	"io"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"
)

func getAvailableFilename(fsPrefix, filename string, mustAppendSuffix bool) string {
	if len(fsPrefix) == 0 {
		fsPrefix = "/"
	} else if fsPrefix[len(fsPrefix)-1] != '/' {
		fsPrefix = fsPrefix + "/"
	}

	if !mustAppendSuffix {
		if _, err := os.Lstat(fsPrefix + filename); os.IsNotExist(err) {
			return filename
		}
	}

	filenamePrefix, filenameSuffix := util.SplitFilename(filename)

	for i := 1; ; i++ {
		newFilename := filenamePrefix + "-" + strconv.Itoa(i) + filenameSuffix
		if _, err := os.Lstat(fsPrefix + newFilename); os.IsNotExist(err) {
			return newFilename
		}
	}

	return ""
}

func (h *handler) saveUploadFiles(fsPrefix string, overwriteExists bool, aliasSubItems []os.FileInfo, r *http.Request) bool {
	errs := []error{}

	reader, err := r.MultipartReader()
	if err != nil {
		errs = append(errs, err)
		return false
	}

	for {
		part, err := reader.NextPart()
		if err != nil {
			if err != io.EOF {
				errs = append(errs, err)
			}
			break
		}

		filename := part.FileName()
		slashIndex := strings.LastIndexAny(filename, "/\\")
		if slashIndex >= 0 {
			filename = filename[slashIndex+1:]
		}
		if len(filename) == 0 {
			continue
		}

		isFilenameAliased := containsItem(aliasSubItems, filename)
		var fsFilename string
		if overwriteExists && !isFilenameAliased {
			tryPath := fsPrefix + "/" + filename
			var info os.FileInfo
			info, err = os.Lstat(tryPath)
			if info != nil && !info.IsDir() {
				err = os.Remove(tryPath)
				if err != nil && !os.IsNotExist(err) {
					errs = append(errs, err)
				}
				// even remove failed, still try to write content to file by TRUNCATE mode
				fsFilename = filename
			}
		}
		if len(fsFilename) == 0 {
			fsFilename = getAvailableFilename(fsPrefix, filename, isFilenameAliased)
		}
		if len(fsFilename) == 0 {
			err := errors.New("no available filename for " + filename)
			errs = append(errs, err)
			continue
		}

		fsPath := path.Clean(fsPrefix + "/" + fsFilename)
		go h.logUpload(filename, fsPath, r)
		file, err := os.OpenFile(fsPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
		if err != nil {
			errs = append(errs, err)
			continue
		}

		_, err = io.Copy(file, part)
		if err != nil {
			errs = append(errs, err)
		}

		err = file.Close()
		if err != nil {
			errs = append(errs, err)
		}
	}

	if len(errs) > 0 {
		go h.logger.LogErrors(errs...)
		return false
	}

	return true
}
