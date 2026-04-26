package report

import (
	"fmt"
	"html/template"
	"strings"

	"github.com/plexusone/pipelineconductor/pkg/model"
)

// CheckHTMLFormatter generates HTML reports for check results.
type CheckHTMLFormatter struct{}

// Format generates an HTML report from check results.
func (f *CheckHTMLFormatter) Format(result *model.CheckResult) ([]byte, error) {
	tmpl, err := template.New("report").Funcs(template.FuncMap{
		"percent": func(n, total int) float64 {
			if total == 0 {
				return 0
			}
			return float64(n) / float64(total) * 100
		},
		"join": strings.Join,
	}).Parse(htmlTemplate)
	if err != nil {
		return nil, fmt.Errorf("parsing template: %w", err)
	}

	// Prepare data
	data := prepareHTMLData(result)

	var buf strings.Builder
	if err := tmpl.Execute(&buf, data); err != nil {
		return nil, fmt.Errorf("executing template: %w", err)
	}

	return []byte(buf.String()), nil
}

// HTMLReportData contains prepared data for the HTML template.
type HTMLReportData struct {
	Result            *model.CheckResult
	FullyCompliant    []model.RepoCheckResult
	Partial           []model.RepoCheckResult
	NonCompliant      []model.RepoCheckResult
	Skipped           []model.RepoCheckResult
	Errors            []model.RepoCheckResult
	ComplianceColor   string
	PartialColor      string
	NonCompliantColor string
}

func prepareHTMLData(result *model.CheckResult) HTMLReportData {
	data := HTMLReportData{
		Result:            result,
		ComplianceColor:   "#22c55e", // green
		PartialColor:      "#eab308", // yellow
		NonCompliantColor: "#ef4444", // red
	}

	for _, repo := range result.Repos {
		if repo.Skipped {
			data.Skipped = append(data.Skipped, repo)
			continue
		}
		if repo.Error != "" {
			data.Errors = append(data.Errors, repo)
			continue
		}

		switch repo.ComplianceLevel {
		case model.ComplianceLevelFull:
			data.FullyCompliant = append(data.FullyCompliant, repo)
		case model.ComplianceLevelPartial:
			data.Partial = append(data.Partial, repo)
		case model.ComplianceLevelNone:
			data.NonCompliant = append(data.NonCompliant, repo)
		}
	}

	return data
}

const htmlTemplate = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Workflow Compliance Report</title>
    <style>
        :root {
            --green: #22c55e;
            --yellow: #eab308;
            --red: #ef4444;
            --blue: #3b82f6;
            --gray-100: #f3f4f6;
            --gray-200: #e5e7eb;
            --gray-700: #374151;
            --gray-900: #111827;
        }

        * {
            box-sizing: border-box;
            margin: 0;
            padding: 0;
        }

        body {
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
            line-height: 1.6;
            color: var(--gray-900);
            background: var(--gray-100);
            padding: 2rem;
        }

        .container {
            max-width: 1200px;
            margin: 0 auto;
        }

        h1 {
            font-size: 2rem;
            margin-bottom: 0.5rem;
        }

        h2 {
            font-size: 1.5rem;
            margin: 2rem 0 1rem;
            padding-bottom: 0.5rem;
            border-bottom: 2px solid var(--gray-200);
        }

        h3 {
            font-size: 1.25rem;
            margin: 1.5rem 0 0.75rem;
        }

        .meta {
            color: var(--gray-700);
            margin-bottom: 2rem;
        }

        .cards {
            display: grid;
            grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
            gap: 1rem;
            margin: 1.5rem 0;
        }

        .card {
            background: white;
            border-radius: 8px;
            padding: 1.5rem;
            box-shadow: 0 1px 3px rgba(0,0,0,0.1);
        }

        .card-title {
            font-size: 0.875rem;
            color: var(--gray-700);
            text-transform: uppercase;
            letter-spacing: 0.05em;
        }

        .card-value {
            font-size: 2.5rem;
            font-weight: bold;
            margin: 0.5rem 0;
        }

        .card-subtitle {
            font-size: 0.875rem;
            color: var(--gray-700);
        }

        .compliant { color: var(--green); }
        .partial { color: var(--yellow); }
        .non-compliant { color: var(--red); }

        .progress-bar {
            height: 24px;
            background: var(--gray-200);
            border-radius: 12px;
            overflow: hidden;
            display: flex;
            margin: 1rem 0;
        }

        .progress-segment {
            height: 100%;
            display: flex;
            align-items: center;
            justify-content: center;
            color: white;
            font-size: 0.75rem;
            font-weight: 600;
            min-width: 30px;
        }

        table {
            width: 100%;
            border-collapse: collapse;
            background: white;
            border-radius: 8px;
            overflow: hidden;
            box-shadow: 0 1px 3px rgba(0,0,0,0.1);
            margin: 1rem 0;
        }

        th, td {
            padding: 0.75rem 1rem;
            text-align: left;
            border-bottom: 1px solid var(--gray-200);
        }

        th {
            background: var(--gray-100);
            font-weight: 600;
            font-size: 0.875rem;
            text-transform: uppercase;
            letter-spacing: 0.05em;
        }

        tr:last-child td {
            border-bottom: none;
        }

        tr:hover {
            background: var(--gray-100);
        }

        .badge {
            display: inline-block;
            padding: 0.25rem 0.5rem;
            border-radius: 4px;
            font-size: 0.75rem;
            font-weight: 600;
        }

        .badge-green { background: #dcfce7; color: #166534; }
        .badge-yellow { background: #fef9c3; color: #854d0e; }
        .badge-red { background: #fee2e2; color: #991b1b; }
        .badge-gray { background: var(--gray-200); color: var(--gray-700); }

        .icon {
            display: inline-block;
            width: 1.25em;
            text-align: center;
        }

        .repo-link {
            color: var(--blue);
            text-decoration: none;
        }

        .repo-link:hover {
            text-decoration: underline;
        }

        .missing-list {
            font-size: 0.875rem;
            color: var(--gray-700);
        }

        .section {
            background: white;
            border-radius: 8px;
            padding: 1.5rem;
            margin: 1.5rem 0;
            box-shadow: 0 1px 3px rgba(0,0,0,0.1);
        }

        .filter-buttons {
            display: flex;
            gap: 0.5rem;
            margin: 1rem 0;
            flex-wrap: wrap;
        }

        .filter-btn {
            padding: 0.5rem 1rem;
            border: 1px solid var(--gray-200);
            background: white;
            border-radius: 4px;
            cursor: pointer;
            font-size: 0.875rem;
        }

        .filter-btn:hover {
            background: var(--gray-100);
        }

        .filter-btn.active {
            background: var(--blue);
            color: white;
            border-color: var(--blue);
        }

        .hidden {
            display: none;
        }

        footer {
            margin-top: 3rem;
            padding-top: 1.5rem;
            border-top: 1px solid var(--gray-200);
            text-align: center;
            color: var(--gray-700);
            font-size: 0.875rem;
        }

        footer a {
            color: var(--blue);
            text-decoration: none;
        }
    </style>
</head>
<body>
    <div class="container">
        <h1>Workflow Compliance Report</h1>
        <p class="meta">
            <strong>Reference:</strong> {{.Result.Config.RefRepo}}@{{.Result.Config.RefBranch}} |
            <strong>Generated:</strong> {{.Result.Timestamp}} |
            <strong>Duration:</strong> {{.Result.ScanDurationMs}}ms
        </p>

        <div class="cards">
            <div class="card">
                <div class="card-title">Total Repositories</div>
                <div class="card-value">{{.Result.Summary.TotalRepos}}</div>
                <div class="card-subtitle">scanned</div>
            </div>
            <div class="card">
                <div class="card-title">Compliance Rate</div>
                <div class="card-value compliant">{{printf "%.1f" .Result.Summary.ComplianceRate}}%</div>
                <div class="card-subtitle">fully compliant</div>
            </div>
            <div class="card">
                <div class="card-title">Fully Compliant</div>
                <div class="card-value compliant">{{.Result.Summary.CompliantRepos}}</div>
                <div class="card-subtitle">repositories</div>
            </div>
            <div class="card">
                <div class="card-title">Partial</div>
                <div class="card-value partial">{{.Result.Summary.PartialRepos}}</div>
                <div class="card-subtitle">repositories</div>
            </div>
            <div class="card">
                <div class="card-title">Non-Compliant</div>
                <div class="card-value non-compliant">{{.Result.Summary.NonCompliant}}</div>
                <div class="card-subtitle">repositories</div>
            </div>
        </div>

        <div class="progress-bar">
            {{if gt .Result.Summary.CompliantRepos 0}}
            <div class="progress-segment" style="background: var(--green); width: {{percent .Result.Summary.CompliantRepos .Result.Summary.TotalRepos}}%">
                {{.Result.Summary.CompliantRepos}}
            </div>
            {{end}}
            {{if gt .Result.Summary.PartialRepos 0}}
            <div class="progress-segment" style="background: var(--yellow); width: {{percent .Result.Summary.PartialRepos .Result.Summary.TotalRepos}}%">
                {{.Result.Summary.PartialRepos}}
            </div>
            {{end}}
            {{if gt .Result.Summary.NonCompliant 0}}
            <div class="progress-segment" style="background: var(--red); width: {{percent .Result.Summary.NonCompliant .Result.Summary.TotalRepos}}%">
                {{.Result.Summary.NonCompliant}}
            </div>
            {{end}}
        </div>

        {{if .Result.Summary.ByLanguage}}
        <h2>By Language</h2>
        <table>
            <thead>
                <tr>
                    <th>Language</th>
                    <th>Total</th>
                    <th>Compliant</th>
                    <th>Rate</th>
                </tr>
            </thead>
            <tbody>
                {{range .Result.Summary.ByLanguage}}
                <tr>
                    <td>{{.Language}}</td>
                    <td>{{.TotalRepos}}</td>
                    <td>{{.CompliantRepos}}</td>
                    <td>{{printf "%.1f" .ComplianceRate}}%</td>
                </tr>
                {{end}}
            </tbody>
        </table>
        {{end}}

        <h2>Repositories</h2>

        <div class="filter-buttons">
            <button class="filter-btn active" data-filter="all">All ({{.Result.Summary.TotalRepos}})</button>
            <button class="filter-btn" data-filter="full">Compliant ({{len .FullyCompliant}})</button>
            <button class="filter-btn" data-filter="partial">Partial ({{len .Partial}})</button>
            <button class="filter-btn" data-filter="none">Non-Compliant ({{len .NonCompliant}})</button>
        </div>

        {{if .NonCompliant}}
        <div class="section repo-section" data-level="none">
            <h3><span class="icon">❌</span> Non-Compliant Repositories</h3>
            <table>
                <thead>
                    <tr>
                        <th>Repository</th>
                        <th>Languages</th>
                        <th>Missing Workflows</th>
                    </tr>
                </thead>
                <tbody>
                    {{range .NonCompliant}}
                    <tr>
                        <td><a href="{{.HTMLURL}}" class="repo-link" target="_blank">{{.FullName}}</a></td>
                        <td>{{join .Languages ", "}}</td>
                        <td class="missing-list">
                            {{range .Missing}}
                            <span class="badge badge-red">{{.WorkflowType}}</span>
                            {{end}}
                        </td>
                    </tr>
                    {{end}}
                </tbody>
            </table>
        </div>
        {{end}}

        {{if .Partial}}
        <div class="section repo-section" data-level="partial">
            <h3><span class="icon">🟡</span> Partially Compliant Repositories</h3>
            <table>
                <thead>
                    <tr>
                        <th>Repository</th>
                        <th>Languages</th>
                        <th>Workflows</th>
                        <th>Issues</th>
                    </tr>
                </thead>
                <tbody>
                    {{range .Partial}}
                    <tr>
                        <td><a href="{{.HTMLURL}}" class="repo-link" target="_blank">{{.FullName}}</a></td>
                        <td>{{join .Languages ", "}}</td>
                        <td>
                            {{range .RequiredWorkflows}}
                            {{if .Present}}
                            <span class="badge {{if eq .MatchType "exact"}}badge-green{{else}}badge-yellow{{end}}">{{.WorkflowType}}</span>
                            {{else}}
                            <span class="badge badge-red">{{.WorkflowType}}</span>
                            {{end}}
                            {{end}}
                        </td>
                        <td>
                            {{range .Missing}}
                            <span class="badge badge-red">Missing: {{.WorkflowType}}</span>
                            {{end}}
                            {{range .RequiredWorkflows}}
                            {{if .FilenameMismatch}}
                            <span class="badge badge-yellow">Filename: {{.WorkflowType}}</span>
                            {{end}}
                            {{end}}
                        </td>
                    </tr>
                    {{end}}
                </tbody>
            </table>
        </div>
        {{end}}

        {{if .FullyCompliant}}
        <div class="section repo-section" data-level="full">
            <h3><span class="icon">✅</span> Fully Compliant Repositories</h3>
            <table>
                <thead>
                    <tr>
                        <th>Repository</th>
                        <th>Languages</th>
                        <th>Workflows</th>
                    </tr>
                </thead>
                <tbody>
                    {{range .FullyCompliant}}
                    <tr>
                        <td><a href="{{.HTMLURL}}" class="repo-link" target="_blank">{{.FullName}}</a></td>
                        <td>{{join .Languages ", "}}</td>
                        <td>
                            {{range .RequiredWorkflows}}
                            <span class="badge badge-green">{{.WorkflowType}}</span>
                            {{end}}
                        </td>
                    </tr>
                    {{end}}
                </tbody>
            </table>
        </div>
        {{end}}

        <h2>Configuration</h2>
        <div class="section">
            <table>
                <tbody>
                    {{if .Result.Config.Orgs}}
                    <tr>
                        <td><strong>Organizations</strong></td>
                        <td>{{join .Result.Config.Orgs ", "}}</td>
                    </tr>
                    {{end}}
                    {{if .Result.Config.Users}}
                    <tr>
                        <td><strong>Users</strong></td>
                        <td>{{join .Result.Config.Users ", "}}</td>
                    </tr>
                    {{end}}
                    <tr>
                        <td><strong>Languages</strong></td>
                        <td>{{join .Result.Config.Languages ", "}}</td>
                    </tr>
                    <tr>
                        <td><strong>Reference Repository</strong></td>
                        <td>{{.Result.Config.RefRepo}}@{{.Result.Config.RefBranch}}</td>
                    </tr>
                    <tr>
                        <td><strong>Strict Mode</strong></td>
                        <td>{{if .Result.Config.Strict}}Yes{{else}}No{{end}}</td>
                    </tr>
                </tbody>
            </table>
        </div>

        <footer>
            Generated by <a href="https://github.com/plexusone/pipelineconductor">PipelineConductor</a>
        </footer>
    </div>

    <script>
        // Filter functionality
        document.querySelectorAll('.filter-btn').forEach(btn => {
            btn.addEventListener('click', () => {
                const filter = btn.dataset.filter;

                // Update active button
                document.querySelectorAll('.filter-btn').forEach(b => b.classList.remove('active'));
                btn.classList.add('active');

                // Filter sections
                document.querySelectorAll('.repo-section').forEach(section => {
                    if (filter === 'all' || section.dataset.level === filter) {
                        section.classList.remove('hidden');
                    } else {
                        section.classList.add('hidden');
                    }
                });
            });
        });
    </script>
</body>
</html>
`
