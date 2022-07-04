package com.emailProject.emailProject.service;

import com.nimbusds.jose.shaded.json.JSONObject;
import io.jsonwebtoken.SignatureAlgorithm;
import io.jsonwebtoken.impl.Base64UrlCodec;
import io.jsonwebtoken.impl.crypto.DefaultJwtSignatureValidator;
import org.springframework.stereotype.Service;

import javax.crypto.spec.SecretKeySpec;

import java.util.Date;
import com.auth0.jwt.JWT;
import com.auth0.jwt.exceptions.JWTDecodeException;
import com.auth0.jwt.interfaces.DecodedJWT;
import static io.jsonwebtoken.SignatureAlgorithm.HS256;
import static java.lang.Long.decode;

@Service
public class AuthorizationService {
    public boolean haveAccess(String token){
        SignatureAlgorithm sa = HS256;
        String secretKey = "secret";
        SecretKeySpec secretKeySpec = new SecretKeySpec(secretKey.getBytes(), sa.getJcaName());
        String[] chunks = token.split("\\.");
        String tokenWithoutSignature = chunks[0] + "." + chunks[1];
        String signature = chunks[2];

        DefaultJwtSignatureValidator validator = new DefaultJwtSignatureValidator(sa, secretKeySpec);

        if (!validator.isValid(tokenWithoutSignature, signature))
            return false;

        DecodedJWT jwt = JWT.decode(token);
        Date date = new Date();
        if( jwt.getExpiresAt().before(date)) {
           return false;
        }
        if(!jwt.getIssuer().equals("hadi"))
            return false;


        return true;
    }
}
