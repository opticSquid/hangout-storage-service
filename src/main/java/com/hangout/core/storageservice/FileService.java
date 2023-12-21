package com.hangout.core.storageservice;

import org.jboss.resteasy.reactive.multipart.FileUpload;

import com.hangout.core.storageservice.dtos.MediaPipelineInit;

import io.vertx.mutiny.core.eventbus.EventBus;
import jakarta.enterprise.context.ApplicationScoped;
import jakarta.inject.Inject;

@ApplicationScoped
public class FileService {
    @Inject
    EventBus bus;

    public String processFile(FileUpload file) {
        String contentType = file.contentType();
        if (!contentType.startsWith("image/") && !contentType.startsWith("video/")) {
            throw new IllegalArgumentException("Invalid file type: " + contentType);
        } else {
            if (contentType.startsWith("image/")) {
                bus.publish("image-process-pipeline-init", new MediaPipelineInit(file.filePath(), file.fileName()));
                return "image uploaded";
            }
            if (contentType.startsWith("video/")) {
                bus.publish("video-process-pipeline-init", new MediaPipelineInit(file.filePath(), file.fileName()));
                return "video uploaded";
            }
            // ? need to verify file checksums before returning something
            return "file uploaded";
        }
    }
}
