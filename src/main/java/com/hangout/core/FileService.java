package com.hangout.core;

import java.io.BufferedReader;
import java.io.IOException;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import io.quarkus.vertx.ConsumeEvent;
import jakarta.enterprise.context.ApplicationScoped;

@ApplicationScoped
public class FileService {
    private static final Logger LOG = LoggerFactory.getLogger(FileService.class);

    @ConsumeEvent(blocking = true, value = "file-service")
    public void processFile(BufferedReader br) throws InterruptedException {

        LOG.info("processFile() begin");

        try (br) {
            String currentLine = null;
            while ((currentLine = br.readLine()) != null) {
                LOG.info("currentLine " + currentLine);
            }
        } catch (IOException e) {
            LOG.error("Error", e);
        }

        LOG.info("processFile() end");

    }
}
