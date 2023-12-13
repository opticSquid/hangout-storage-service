package com.hangout.core;

import java.io.BufferedReader;
import java.io.IOException;
import java.nio.file.Files;
import java.util.List;

import org.eclipse.microprofile.openapi.annotations.enums.SchemaType;
import org.eclipse.microprofile.openapi.annotations.media.Schema;
import org.eclipse.microprofile.openapi.annotations.responses.APIResponse;
import org.jboss.resteasy.reactive.RestForm;
import org.jboss.resteasy.reactive.multipart.FileUpload;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import io.vertx.mutiny.core.eventbus.EventBus;
import jakarta.enterprise.context.RequestScoped;
import jakarta.inject.Inject;
import jakarta.ws.rs.Consumes;
import jakarta.ws.rs.POST;
import jakarta.ws.rs.Path;
import jakarta.ws.rs.core.MediaType;
import jakarta.ws.rs.core.Response;

@RequestScoped
@Path("upload")
public class UploadResource {
    private static final Logger LOG = LoggerFactory.getLogger(UploadResource.class);

    @Inject
    EventBus bus;

    @POST
    @Consumes(MediaType.MULTIPART_FORM_DATA)
    @APIResponse(responseCode = "202")
    public Response upload(MultipartBody body) throws IOException {

        LOG.info("upload() quantity of files + " + body.files.size());

        for (FileUpload file : body.files) {

            LOG.info("filePath " + file.filePath());

            BufferedReader br = Files.newBufferedReader(file.filePath());

            bus.send("file-service", br);
        }

        LOG.info("upload() before response Accepted");

        return Response
                .accepted()
                .build();
    }

    // Class that will define the OpenAPI schema for the binary type input (upload)
    @Schema(type = SchemaType.STRING, format = "binary")
    public interface UploadItemSchema {
    }

    // We instruct OpenAPI to use the schema provided by the 'UploadFormSchema'
    // class implementation and thus define a valid OpenAPI schema for the Swagger
    // UI
    public static class MultipartBody {
        @Schema(implementation = UploadItemSchema[].class)
        @RestForm("files")
        public List<FileUpload> files;
    }

}
