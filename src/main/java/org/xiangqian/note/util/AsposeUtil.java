package org.xiangqian.note.util;

import com.aspose.cells.Workbook;
import com.aspose.words.Document;
import org.apache.commons.io.FileUtils;
import org.apache.commons.io.IOUtils;

import java.io.*;
import java.nio.charset.StandardCharsets;

/**
 * https://docs.aspose.com
 *
 * @author xiangqian
 * @date 21:08 2024/04/28
 */
public class AsposeUtil {

    /**
     * xls转html
     *
     * @param xlsFile
     * @param htmlFile
     * @param attDir   附件目录
     * @throws Exception
     */
    public static void convertXlsToHtml(File xlsFile, File htmlFile, File attDir) throws Exception {
        InputStream inputStream = null;
        OutputStream outputStream = null;
        try {
            inputStream = new FileInputStream(xlsFile);
            Workbook workbook = new Workbook(inputStream);
            outputStream = new FileOutputStream(htmlFile);
            com.aspose.cells.HtmlSaveOptions saveOptions = new com.aspose.cells.HtmlSaveOptions();
            if (attDir != null) {
                saveOptions.setAttachedFilesDirectory(attDir.getAbsolutePath());
            }
            workbook.save(outputStream, saveOptions);
        } finally {
            IOUtils.closeQuietly(outputStream, inputStream);
        }
        System.out.println(FileUtils.readFileToString(htmlFile, StandardCharsets.UTF_8));
    }

    /**
     * doc转pdf
     *
     * @param docFile
     * @param pdfFile
     * @throws Exception
     */
    public static void convertDocToPdf(File docFile, File pdfFile) throws Exception {
        InputStream inputStream = null;
        OutputStream outputStream = null;
        try {
            // 支持 DOC, DOCX, OOXML, RTF HTML, OpenDocument, PDF, EPUB, XPS, SWF 相互转换
            inputStream = new FileInputStream(docFile);
            Document document = new Document(inputStream);
            outputStream = new FileOutputStream(pdfFile);
            document.save(outputStream, com.aspose.words.SaveFormat.PDF);
        } finally {
            IOUtils.closeQuietly(outputStream, inputStream);
        }
    }

}
