package gui

// indexTemplate is the embedded HTML template for the web GUI
const indexTemplate = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>CodexGigantus - Configuration GUI</title>
    <style>
        * {
            margin: 0;
            padding: 0;
            box-sizing: border-box;
        }
        body {
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Oxygen, Ubuntu, Cantarell, sans-serif;
            line-height: 1.6;
            color: #333;
            background: #f5f5f5;
        }
        .container {
            max-width: 1200px;
            margin: 0 auto;
            padding: 20px;
        }
        header {
            background: #2c3e50;
            color: white;
            padding: 20px 0;
            margin-bottom: 30px;
        }
        h1 {
            text-align: center;
            font-size: 2em;
        }
        .card {
            background: white;
            border-radius: 8px;
            padding: 25px;
            margin-bottom: 20px;
            box-shadow: 0 2px 4px rgba(0,0,0,0.1);
        }
        .card h2 {
            margin-bottom: 20px;
            color: #2c3e50;
            border-bottom: 2px solid #3498db;
            padding-bottom: 10px;
        }
        .form-group {
            margin-bottom: 20px;
        }
        label {
            display: block;
            margin-bottom: 5px;
            font-weight: 600;
            color: #555;
        }
        input[type="text"],
        input[type="number"],
        input[type="password"],
        select,
        textarea {
            width: 100%;
            padding: 10px;
            border: 1px solid #ddd;
            border-radius: 4px;
            font-size: 14px;
        }
        input[type="checkbox"] {
            margin-right: 5px;
        }
        .checkbox-group {
            display: flex;
            align-items: center;
        }
        button {
            background: #3498db;
            color: white;
            padding: 12px 24px;
            border: none;
            border-radius: 4px;
            cursor: pointer;
            font-size: 16px;
            margin-right: 10px;
        }
        button:hover {
            background: #2980b9;
        }
        button.secondary {
            background: #95a5a6;
        }
        button.secondary:hover {
            background: #7f8c8d;
        }
        button.success {
            background: #27ae60;
        }
        button.success:hover {
            background: #229954;
        }
        .tabs {
            display: flex;
            border-bottom: 2px solid #ddd;
            margin-bottom: 20px;
        }
        .tab {
            padding: 10px 20px;
            cursor: pointer;
            border: none;
            background: none;
            margin: 0;
        }
        .tab.active {
            border-bottom: 3px solid #3498db;
            color: #3498db;
            font-weight: 600;
        }
        .tab-content {
            display: none;
        }
        .tab-content.active {
            display: block;
        }
        .alert {
            padding: 15px;
            margin-bottom: 20px;
            border-radius: 4px;
        }
        .alert-success {
            background: #d4edda;
            border: 1px solid #c3e6cb;
            color: #155724;
        }
        .alert-error {
            background: #f8d7da;
            border: 1px solid #f5c6cb;
            color: #721c24;
        }
        .hidden {
            display: none;
        }
        #output {
            max-height: 400px;
            overflow-y: auto;
            background: #f8f9fa;
            padding: 15px;
            border-radius: 4px;
            font-family: monospace;
            font-size: 12px;
            white-space: pre-wrap;
        }
        .row {
            display: grid;
            grid-template-columns: 1fr 1fr;
            gap: 15px;
        }
        @media (max-width: 768px) {
            .row {
                grid-template-columns: 1fr;
            }
        }
    </style>
</head>
<body>
    <header>
        <div class="container">
            <h1>CodexGigantus Configuration</h1>
        </div>
    </header>

    <div class="container">
        <div id="message"></div>

        <div class="card">
            <h2>Configuration Management</h2>
            <div class="form-group">
                <label>Config Name:</label>
                <input type="text" id="configName" placeholder="My Configuration">
            </div>
            <div class="form-group">
                <label>Description:</label>
                <input type="text" id="configDesc" placeholder="Configuration description">
            </div>
            <div>
                <button onclick="saveConfig()">💾 Save Config</button>
                <button class="secondary" onclick="loadConfig()">📂 Load Config</button>
            </div>
        </div>

        <div class="card">
            <h2>Source Configuration</h2>

            <div class="tabs">
                <button class="tab active" onclick="switchTab('filesystem')">Filesystem</button>
                <button class="tab" onclick="switchTab('csv')">CSV/TSV</button>
                <button class="tab" onclick="switchTab('database')">Database</button>
            </div>

            <!-- Filesystem Tab -->
            <div id="filesystem-tab" class="tab-content active">
                <div class="form-group">
                    <label>Directories (comma-separated):</label>
                    <input type="text" id="directories" value="." placeholder="./src,./pkg">
                </div>
                <div class="form-group checkbox-group">
                    <input type="checkbox" id="recursive" checked>
                    <label for="recursive">Recursive</label>
                </div>
                <div class="form-group">
                    <label>Ignore Files (comma-separated):</label>
                    <input type="text" id="ignoreFiles" placeholder=".DS_Store,*.log">
                </div>
                <div class="form-group">
                    <label>Ignore Directories (comma-separated):</label>
                    <input type="text" id="ignoreDirs" placeholder="node_modules,.git,vendor">
                </div>
                <div class="form-group">
                    <label>Exclude Extensions (comma-separated, without dot):</label>
                    <input type="text" id="excludeExts" placeholder="log,tmp">
                </div>
                <div class="form-group">
                    <label>Include Extensions Only (comma-separated, without dot):</label>
                    <input type="text" id="includeExts" placeholder="go,py,js">
                </div>
            </div>

            <!-- CSV Tab -->
            <div id="csv-tab" class="tab-content">
                <div class="form-group">
                    <label>CSV/TSV File Path:</label>
                    <input type="text" id="csvPath" placeholder="/path/to/data.csv">
                </div>
                <div class="row">
                    <div class="form-group">
                        <label>Delimiter:</label>
                        <select id="csvDelimiter">
                            <option value=",">Comma (,)</option>
                            <option value="\t">Tab (\t)</option>
                        </select>
                    </div>
                    <div class="form-group checkbox-group">
                        <input type="checkbox" id="csvHeader" checked>
                        <label for="csvHeader">Has Header Row</label>
                    </div>
                </div>
                <div class="row">
                    <div class="form-group">
                        <label>Path Column Index:</label>
                        <input type="number" id="csvPathCol" value="0" min="0">
                    </div>
                    <div class="form-group">
                        <label>Content Column Index:</label>
                        <input type="number" id="csvContentCol" value="1" min="0">
                    </div>
                </div>
            </div>

            <!-- Database Tab -->
            <div id="database-tab" class="tab-content">
                <div class="row">
                    <div class="form-group">
                        <label>Database Type:</label>
                        <select id="dbType">
                            <option value="postgres">PostgreSQL</option>
                            <option value="mysql">MySQL</option>
                            <option value="sqlite">SQLite</option>
                        </select>
                    </div>
                    <div class="form-group">
                        <label>Host:</label>
                        <input type="text" id="dbHost" value="localhost">
                    </div>
                </div>
                <div class="row">
                    <div class="form-group">
                        <label>Port:</label>
                        <input type="number" id="dbPort" value="5432">
                    </div>
                    <div class="form-group">
                        <label>Database Name:</label>
                        <input type="text" id="dbName" placeholder="codex">
                    </div>
                </div>
                <div class="row">
                    <div class="form-group">
                        <label>User:</label>
                        <input type="text" id="dbUser" placeholder="postgres">
                    </div>
                    <div class="form-group">
                        <label>Password:</label>
                        <input type="password" id="dbPassword">
                    </div>
                </div>
                <div class="form-group">
                    <label>Table Name:</label>
                    <input type="text" id="dbTable" value="code_files">
                </div>
                <div class="row">
                    <div class="form-group">
                        <label>Path Column:</label>
                        <input type="text" id="dbColPath" value="file_path">
                    </div>
                    <div class="form-group">
                        <label>Content Column:</label>
                        <input type="text" id="dbColContent" value="content">
                    </div>
                </div>
                <div class="form-group">
                    <label>Custom Query (optional):</label>
                    <textarea id="dbQuery" rows="3" placeholder="SELECT file_path, content FROM code_files WHERE..."></textarea>
                </div>
                <button class="secondary" onclick="testDatabase()">🔌 Test Connection</button>
            </div>
        </div>

        <div class="card">
            <h2>Output Configuration</h2>
            <div class="form-group">
                <label>Output File:</label>
                <input type="text" id="outputFile" value="output.txt">
            </div>
            <div class="form-group checkbox-group">
                <input type="checkbox" id="showSize">
                <label for="showSize">Show Size</label>
            </div>
            <div class="form-group checkbox-group">
                <input type="checkbox" id="showFuncs">
                <label for="showFuncs">Show Function Signatures Only (Go files)</label>
            </div>
            <div class="form-group checkbox-group">
                <input type="checkbox" id="debug">
                <label for="debug">Debug Mode</label>
            </div>
        </div>

        <div class="card">
            <h2>Process Files</h2>
            <button class="success" onclick="processFiles()">▶️ Process Files</button>
            <div id="output" class="hidden"></div>
        </div>
    </div>

    <script>
        let currentSourceType = 'filesystem';

        function switchTab(tabName) {
            currentSourceType = tabName;

            // Update tabs
            document.querySelectorAll('.tab').forEach(tab => {
                tab.classList.remove('active');
            });
            event.target.classList.add('active');

            // Update content
            document.querySelectorAll('.tab-content').forEach(content => {
                content.classList.remove('active');
            });
            document.getElementById(tabName + '-tab').classList.add('active');
        }

        function showMessage(message, type) {
            const messageDiv = document.getElementById('message');
            messageDiv.innerHTML = '<div class="alert alert-' + type + '">' + message + '</div>';
            setTimeout(() => {
                messageDiv.innerHTML = '';
            }, 5000);
        }

        function getConfig() {
            const config = {
                source_type: currentSourceType,
                name: document.getElementById('configName').value,
                description: document.getElementById('configDesc').value,
                output_file: document.getElementById('outputFile').value,
                show_size: document.getElementById('showSize').checked,
                show_funcs: document.getElementById('showFuncs').checked,
                debug: document.getElementById('debug').checked
            };

            if (currentSourceType === 'filesystem') {
                config.directories = document.getElementById('directories').value.split(',').map(s => s.trim());
                config.recursive = document.getElementById('recursive').checked;
                config.ignore_files = document.getElementById('ignoreFiles').value.split(',').map(s => s.trim()).filter(s => s);
                config.ignore_dirs = document.getElementById('ignoreDirs').value.split(',').map(s => s.trim()).filter(s => s);
                config.exclude_extensions = document.getElementById('excludeExts').value.split(',').map(s => s.trim()).filter(s => s);
                config.include_extensions = document.getElementById('includeExts').value.split(',').map(s => s.trim()).filter(s => s);
            } else if (currentSourceType === 'csv' || currentSourceType === 'tsv') {
                config.csv_file_path = document.getElementById('csvPath').value;
                config.csv_delimiter = document.getElementById('csvDelimiter').value;
                config.csv_path_column = parseInt(document.getElementById('csvPathCol').value);
                config.csv_content_column = parseInt(document.getElementById('csvContentCol').value);
                config.csv_has_header = document.getElementById('csvHeader').checked;
            } else if (currentSourceType === 'database') {
                config.db_type = document.getElementById('dbType').value;
                config.db_host = document.getElementById('dbHost').value;
                config.db_port = parseInt(document.getElementById('dbPort').value);
                config.db_name = document.getElementById('dbName').value;
                config.db_user = document.getElementById('dbUser').value;
                config.db_password = document.getElementById('dbPassword').value;
                config.db_table_name = document.getElementById('dbTable').value;
                config.db_column_path = document.getElementById('dbColPath').value;
                config.db_column_content = document.getElementById('dbColContent').value;
                config.db_query = document.getElementById('dbQuery').value;
            }

            return config;
        }

        async function saveConfig() {
            const filename = prompt('Enter filename to save (e.g., config.json or config.yaml):');
            if (!filename) return;

            const config = getConfig();

            try {
                const response = await fetch('/api/config/save', {
                    method: 'POST',
                    headers: {'Content-Type': 'application/json'},
                    body: JSON.stringify({file_path: filename})
                });

                if (response.ok) {
                    showMessage('Configuration saved to ' + filename, 'success');
                } else {
                    const error = await response.text();
                    showMessage('Error: ' + error, 'error');
                }
            } catch (error) {
                showMessage('Error: ' + error.message, 'error');
            }
        }

        async function loadConfig() {
            const filename = prompt('Enter filename to load (e.g., config.json or config.yaml):');
            if (!filename) return;

            try {
                const response = await fetch('/api/config/load', {
                    method: 'POST',
                    headers: {'Content-Type': 'application/json'},
                    body: JSON.stringify({file_path: filename})
                });

                if (response.ok) {
                    const data = await response.json();
                    showMessage('Configuration loaded from ' + filename, 'success');
                    location.reload(); // Reload to update UI with new config
                } else {
                    const error = await response.text();
                    showMessage('Error: ' + error, 'error');
                }
            } catch (error) {
                showMessage('Error: ' + error.message, 'error');
            }
        }

        async function testDatabase() {
            const config = getConfig();

            try {
                const response = await fetch('/api/config', {
                    method: 'POST',
                    headers: {'Content-Type': 'application/json'},
                    body: JSON.stringify(config)
                });

                if (!response.ok) {
                    const error = await response.text();
                    showMessage('Config Error: ' + error, 'error');
                    return;
                }

                const testResponse = await fetch('/api/test-db', {
                    method: 'POST'
                });

                if (testResponse.ok) {
                    showMessage('Database connection successful!', 'success');
                } else {
                    const error = await testResponse.text();
                    showMessage('Connection Error: ' + error, 'error');
                }
            } catch (error) {
                showMessage('Error: ' + error.message, 'error');
            }
        }

        async function processFiles() {
            const config = getConfig();

            try {
                // Update config first
                const configResponse = await fetch('/api/config', {
                    method: 'POST',
                    headers: {'Content-Type': 'application/json'},
                    body: JSON.stringify(config)
                });

                if (!configResponse.ok) {
                    const error = await configResponse.text();
                    showMessage('Config Error: ' + error, 'error');
                    return;
                }

                // Process files
                showMessage('Processing files...', 'success');
                const response = await fetch('/api/process', {
                    method: 'POST'
                });

                if (response.ok) {
                    const data = await response.json();
                    const outputDiv = document.getElementById('output');
                    outputDiv.classList.remove('hidden');
                    outputDiv.textContent =
                        'Processed ' + data.file_count + ' files\n' +
                        'Output size: ' + data.output_size + ' bytes\n' +
                        'Saved to: ' + data.output_file + '\n\n' +
                        'Preview:\n' + data.output.substring(0, 5000) +
                        (data.output.length > 5000 ? '\n\n... (truncated)' : '');
                    showMessage('Processing complete!', 'success');
                } else {
                    const error = await response.text();
                    showMessage('Processing Error: ' + error, 'error');
                }
            } catch (error) {
                showMessage('Error: ' + error.message, 'error');
            }
        }
    </script>
</body>
</html>
`
