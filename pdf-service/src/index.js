const express = require('express');
const cors = require('cors');
const { generatePDF } = require('./pdf');

const app = express();
const PORT = process.env.PORT || 3001;

// Increase body limit for large HTML strings
app.use(express.json({ limit: '50mb' }));
app.use(express.urlencoded({ limit: '50mb', extended: true }));
app.use(cors());

/**
 * POST /generate-pdf
 * Expects { html: string }
 * Returns application/pdf binary
 */
app.post('/generate-pdf', async (req, res) => {
    const { html } = req.body;

    if (!html) {
        return res.status(400).json({ error: 'HTML content is required.' });
    }

    try {
        console.log(`[${new Date().toISOString()}] Generating PDF...`);
        const pdfBuffer = await generatePDF(html);

        res.setHeader('Content-Type', 'application/pdf');
        res.setHeader('Content-Disposition', 'attachment; filename=case-report.pdf');
        res.send(pdfBuffer);
        console.log(`[${new Date().toISOString()}] PDF generated successfully.`);
    } catch (error) {
        console.error('PDF Generation Error:', error);
        res.status(500).json({
            error: 'Failed to generate PDF.',
            details: error.message,
        });
    }
});

// Health check endpoint
app.get('/health', (req, res) => {
    res.json({ status: 'ok', service: 'pdf-generation' });
});

app.listen(PORT, '0.0.0.0', () => {
    console.log(`PDF Generation Microservice listening on port ${PORT}`);
});
