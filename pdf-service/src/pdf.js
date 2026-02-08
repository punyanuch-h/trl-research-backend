const puppeteer = require('puppeteer');
const path = require('path');
const fs = require('fs');

/**
 * Generates a PDF from HTML content using Puppeteer.
 * @param {string} htmlContent - Full HTML string to render.
 * @returns {Promise<Buffer>} - PDF binary buffer.
 */
async function generatePDF(htmlContent) {
    let browser;
    try {
        // Launch puppeteer with recommended flags for server environments
        browser = await puppeteer.launch({
            headless: 'new',
            args: [
                '--no-sandbox',
                '--disable-setuid-sandbox',
                '--disable-dev-shm-usage',
                '--font-render-hinting=none', // Better font rendering for PDFs
            ],
        });

        const page = await browser.newPage();

        // Set A4 viewport for consistent rendering calculation
        await page.setViewport({
            width: 794, // 210mm at 96 DPI
            height: 1123, // 297mm at 96 DPI
            deviceScaleFactor: 2,
        });

        // Prepare fonts to be injected as base64
        const fontPath = path.resolve(__dirname, '../assets/fonts');
        const fontsToLoad = [
            { file: 'Sarabun-Regular.ttf', weight: 'normal', style: 'normal' },
            { file: 'Sarabun-Bold.ttf', weight: 'bold', style: 'normal' },
            { file: 'Sarabun-Italic.ttf', weight: 'normal', style: 'italic' }
        ];

        let fontStyles = '';
        for (const font of fontsToLoad) {
            const fullPath = path.join(fontPath, font.file);
            if (fs.existsSync(fullPath)) {
                const base64Font = fs.readFileSync(fullPath).toString('base64');
                fontStyles += `
          @font-face {
            font-family: 'Sarabun';
            src: url(data:font/ttf;charset=utf-8;base64,${base64Font}) format('truetype');
            font-weight: ${font.weight};
            font-style: ${font.style};
          }
        `;
            }
        }

        // Global styles for PDF layout and Thai rendering
        const pdfStyles = `
      ${fontStyles}
      body {
        font-family: 'Sarabun', sans-serif !important;
        line-height: 1.5;
        color: #000;
        margin: 0;
        padding: 0;
      }
      /* Print specific optimizations */
      @media print {
        body {
          -webkit-print-color-adjust: exact;
        }
      }
      /* Helper class for manual page breaks */
      .page-break {
        page-break-before: always;
      }
      /* Thai language support: Rely on Chrome's layout engine */
      * {
        word-wrap: break-word;
      }
    `;

        // Set the HTML content
        // We ensure the HTML has the lang="th" attribute to assist Chrome's line-breaking engine
        let fullHtml = htmlContent;
        if (!fullHtml.includes('<html')) {
            fullHtml = `<html lang="th"><body>${htmlContent}</body></html>`;
        } else if (!fullHtml.includes('lang=')) {
            fullHtml = fullHtml.replace('<html', '<html lang="th"');
        }

        await page.setContent(fullHtml, {
            waitUntil: ['load', 'networkidle0'],
        });

        // Inject font styles and global PDF overrides
        await page.addStyleTag({ content: pdfStyles });

        // Generate PDF buffer
        const pdfBuffer = await page.pdf({
            format: 'A4',
            printBackground: true,
            margin: {
                top: '20mm',
                right: '20mm',
                bottom: '20mm',
                left: '20mm',
            },
            displayHeaderFooter: false,
            preferCSSPageSize: true, // Respect @page CSS if present
        });

        return pdfBuffer;
    } catch (error) {
        console.error('Puppeteer Error:', error);
        throw error;
    } finally {
        if (browser) {
            await browser.close();
        }
    }
}

module.exports = { generatePDF };
