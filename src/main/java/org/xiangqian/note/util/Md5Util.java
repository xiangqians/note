package org.xiangqian.note.util;

import org.apache.commons.codec.digest.DigestUtils;

import java.io.IOException;
import java.io.InputStream;

/**
 * @author xiangqian
 * @date 20:36 2023/11/14
 */
public class Md5Util {

    public static byte[] encrypt(InputStream data) throws IOException {
        return DigestUtils.md5(data);
    }

    public static byte[] encrypt(byte[] data) {
        return DigestUtils.md5(data);
    }

    public static byte[] encrypt(String data) {
        return DigestUtils.md5(data);
    }

    public static String encryptHex(InputStream data) throws IOException {
        return DigestUtils.md5Hex(data);
    }

    public static String encryptHex(byte[] data) {
        return DigestUtils.md5Hex(data);
    }

    public static String encryptHex(String data) {
        return DigestUtils.md5Hex(data);
    }

}
