package com.hangout.core.storageservice.media_processors;

import java.awt.image.BufferedImage;
import java.io.File;
import java.io.FileOutputStream;
import java.io.IOException;
import java.io.OutputStream;
import java.util.Iterator;
import java.util.UUID;

import javax.imageio.IIOImage;
import javax.imageio.ImageIO;
import javax.imageio.ImageWriteParam;
import javax.imageio.ImageWriter;
import javax.imageio.stream.ImageOutputStream;

import org.apache.commons.io.FileUtils;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import com.hangout.core.storageservice.dtos.MediaPipelineInit;

import io.quarkus.vertx.ConsumeEvent;
import io.smallrye.mutiny.Uni;
import jakarta.enterprise.context.ApplicationScoped;
import jakarta.ws.rs.ProcessingException;

@ApplicationScoped
public class ImageProcessor {
    private static final Logger LOG = LoggerFactory.getLogger(ImageProcessor.class);

    @ConsumeEvent("image-process-pipeline-init")
    public Uni<Void> compressImage(MediaPipelineInit media) {
        LOG.info("file path: {}", media.filePath().toString());
        try {
            File imageFile = media.filePath().toFile();
            File destFile = new File("store", UUID.randomUUID().toString() + ".jpg");
            if (FileUtils.sizeOf(imageFile) <= 2_000_000L) {
                LOG.debug("file is less than 2 MB, size:{}", FileUtils.sizeOf(imageFile));
                FileUtils.copyFile(imageFile, destFile);
            } else {
                LOG.debug("file is larger than 2 MB, size:{}", FileUtils.sizeOf(imageFile));
                BufferedImage image = ImageIO.read(imageFile);
                OutputStream os = new FileOutputStream(destFile);
                Iterator<ImageWriter> writers = ImageIO.getImageWritersByFormatName("jpg");
                ImageWriter writer = (ImageWriter) writers.next();

                ImageOutputStream ios = ImageIO.createImageOutputStream(os);
                writer.setOutput(ios);

                ImageWriteParam param = writer.getDefaultWriteParam();

                param.setCompressionMode(ImageWriteParam.MODE_EXPLICIT);
                // compress to 60% quality
                param.setCompressionQuality(0.60f); // Change the quality value you prefer
                writer.write(null, new IIOImage(image, null, null), param);
                os.close();
                ios.close();
                writer.dispose();
            }
            imageFile.delete();
            return Uni.createFrom().voidItem();
        } catch (IOException e) {
            throw new ProcessingException("Image file could not be processed");
        }
    }
}
