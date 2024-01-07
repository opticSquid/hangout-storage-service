package com.hangout.core.storageservice;

import com.hangout.core.storageservice.dtos.FileUploadResponseDTO;
import org.jboss.resteasy.reactive.multipart.FileUpload;

import com.hangout.core.storageservice.dtos.MediaPipelineInit;

import io.vertx.mutiny.core.eventbus.EventBus;
import jakarta.enterprise.context.ApplicationScoped;
import jakarta.inject.Inject;

import java.io.File;

@ApplicationScoped
public class FileService {
    @Inject
    EventBus bus;

    public FileUploadResponseDTO processFile(FileUpload file) {
        String contentType = file.contentType();
        if (!contentType.startsWith("image/") && !contentType.startsWith("video/")) {
            throw new IllegalArgumentException("Invalid file type: " + contentType);
        } else {
            if (contentType.startsWith("image/")) {
                bus.publish("image-process-pipeline-init", new MediaPipelineInit(file.filePath(), file.fileName()));
                return new FileUploadResponseDTO("image uploaded");
            }
            if (contentType.startsWith("video/")) {
                bus.publish("video-process-pipeline-init", new MediaPipelineInit(file.filePath(), file.fileName()));
                return new FileUploadResponseDTO("video uploaded");
            }
            // ? need to verify file checksums before returning something
            return new FileUploadResponseDTO("file uploaded");
        }
    }
}
