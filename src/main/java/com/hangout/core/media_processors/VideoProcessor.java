package com.hangout.core.media_processors;

import java.io.File;
import java.io.NotActiveException;
import java.util.UUID;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import com.hangout.core.dtos.MediaPipelineInit;

import io.quarkus.vertx.ConsumeEvent;
import io.smallrye.mutiny.Uni;
import jakarta.ws.rs.ProcessingException;

public class VideoProcessor {
    private static final Logger LOG = LoggerFactory.getLogger(ImageProcessor.class);

    @ConsumeEvent("video-process-pipeline-init")
    public Uni<Void> compressVideo(MediaPipelineInit media) {
        LOG.info("file path: {}", media.filePath().toString());
        try {
            File videoFile = media.filePath().toFile();
            File destFile = new File("store", UUID.randomUUID().toString() + ".mp4");
            return Uni.createFrom().voidItem();

        } catch (Exception ex) {
            throw new ProcessingException("This feature is not active yet");
        }
    }
}
