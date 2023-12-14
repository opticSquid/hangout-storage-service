package com.hangout.core;

import java.io.IOException;

import org.eclipse.microprofile.openapi.annotations.enums.SchemaType;
import org.eclipse.microprofile.openapi.annotations.media.Schema;
import org.jboss.resteasy.reactive.RestForm;
import org.jboss.resteasy.reactive.multipart.FileUpload;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import io.smallrye.mutiny.Uni;
import io.vertx.mutiny.core.eventbus.EventBus;
import jakarta.enterprise.context.RequestScoped;
import jakarta.inject.Inject;
import jakarta.ws.rs.Consumes;
import jakarta.ws.rs.POST;
import jakarta.ws.rs.Path;
import jakarta.ws.rs.Produces;
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
    @Produces(MediaType.APPLICATION_JSON)
    public Uni<Response> upload(MultipartBody body) throws IOException {
        bus.publish("file-service", body.file);
        bus.<Uni<String>>consumer("file-path", message -> LOG.info("file path is: {}", message.body()));
        return Uni.createFrom().item(Response.ok().build());
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
        @RestForm("file")
        public FileUpload file;
        public FileType fileType;
    }

}
