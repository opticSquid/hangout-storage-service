package com.hangout.core;

import java.nio.file.Path;

import org.apache.commons.io.FileUtils;
import org.jboss.resteasy.reactive.multipart.FileUpload;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import io.quarkus.vertx.ConsumeEvent;
import io.vertx.mutiny.core.eventbus.EventBus;
import jakarta.enterprise.context.ApplicationScoped;
import jakarta.inject.Inject;

@ApplicationScoped
public class FileService {
    private static final Logger LOG = LoggerFactory.getLogger(FileService.class);
    private String storagePath = "./store";

    @Inject
    EventBus bus;

    @ConsumeEvent(blocking = true, value = "file-service")
    public void processFile(FileUpload file) {
        String fileName = file.fileName();
        String contentType = file.contentType();
        if (!contentType.startsWith("image/")) {
            throw new IllegalArgumentException("Invalid file type: " + contentType);
        }
        // String uniqueFileName = generateUniqueFileName(fileName);
        // String targetPath = storagePath + uniqueFileName;
        Path filePath = file.filePath();

        bus.publish("file-path", filePath.toString());
    }

    private String generateUniqueFileName(String fileName) {
        return null;
    }
}
