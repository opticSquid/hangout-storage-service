package com.hangout.core.dtos;

import org.jboss.resteasy.reactive.RestForm;
import org.jboss.resteasy.reactive.multipart.FileUpload;

public class IncomingRequestBody {
    @RestForm("file")
    public FileUpload file;
}
