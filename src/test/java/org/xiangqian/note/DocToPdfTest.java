package org.xiangqian.note;

import com.aspose.words.Document;
import com.aspose.words.SaveFormat;

/**
 * @author xiangqian
 * @date 17:07 2023/07/10
 */
public class DocToPdfTest {

    public static void main(String[] args) throws Exception {
        String input = "E:\\tmp\\convert\\test.docx";
        String output = "E:\\tmp\\convert\\test.pdf";
        Document doc = new Document(input);
        // 全面支持DOC, DOCX, OOXML, RTF HTML, OpenDocument, PDF, EPUB, XPS, SWF 相互转换
        doc.save(output, SaveFormat.PDF);
    }

}
