package com.hangout.core;

import org.jboss.logging.Logger;
import org.jboss.resteasy.reactive.multipart.FileUpload;

import jakarta.enterprise.context.ApplicationScoped;

@ApplicationScoped
public class FileService {
    private static final Logger LOG = Logger.getLogger(FileService.class);

    // @ConsumeEvent(blocking = true, value = "file-service")
    public String processFile(FileUpload file) {
        String fileName = file.fileName();
        String contentType = file.contentType();
        if (!contentType.startsWith("image/") && !contentType.startsWith("video/")) {
            throw new IllegalArgumentException("Invalid file type: " + contentType);
        }
        return file.filePath().toString() + "/" + fileName;
    }
}
