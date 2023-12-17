package com.hangout.core;

import com.hangout.core.dtos.IncomingRequestBody;

import jakarta.enterprise.context.RequestScoped;
import jakarta.inject.Inject;
import jakarta.ws.rs.Consumes;
import jakarta.ws.rs.POST;
import jakarta.ws.rs.Path;
import jakarta.ws.rs.Produces;
import jakarta.ws.rs.core.MediaType;

@Path("upload")
@RequestScoped
public class UploadResource {
    @Inject
    FileService fs;

    @POST
    @Consumes(MediaType.MULTIPART_FORM_DATA)
    @Produces(MediaType.APPLICATION_JSON)
    public String upload(IncomingRequestBody body) {
        return fs.processFile(body.file);
    }

    // Class that will define the OpenAPI schema for the binary type input (upload)
    // @Schema(type = SchemaType.STRING, format = "binary")
    // public interface UploadItemSchema {
    // }

    // We instruct OpenAPI to use the schema provided by the 'UploadFormSchema'
    // class implementation and thus define a valid OpenAPI schema for the Swagger
    // UI
    // public static class MultipartBody {
    // @Schema(implementation = UploadItemSchema[].class)
    // @RestForm("file")
    // public FileUpload file;
    // }

}
