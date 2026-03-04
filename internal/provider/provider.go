package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"google.golang.org/api/transport"
)

var _ provider.Provider = &FirestoreProvider{}

type FirestoreProvider struct {
	version string
}

type FirestoreProviderModel struct {
	Project     types.String `tfsdk:"project"`
	Credentials types.String `tfsdk:"credentials"`
	Database    types.String `tfsdk:"database"`
}

type FirestoreClient struct {
	HTTPClient *http.Client
	Project    string
	Database   string
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &FirestoreProvider{
			version: version,
		}
	}
}

func (p *FirestoreProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "firestore"
	resp.Version = p.version
}

func (p *FirestoreProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Interact with Google Cloud Firestore.",
		Attributes: map[string]schema.Attribute{
			"project": schema.StringAttribute{
				Description: "The GCP project ID. Can also be set via GOOGLE_PROJECT or GOOGLE_CLOUD_PROJECT environment variables.",
				Optional:    true,
			},
			"credentials": schema.StringAttribute{
				Description: "Service account JSON credentials. Can also be set via GOOGLE_CREDENTIALS or GOOGLE_APPLICATION_CREDENTIALS environment variables.",
				Optional:    true,
				Sensitive:   true,
			},
			"database": schema.StringAttribute{
				Description: "The Firestore database ID. Defaults to '(default)'.",
				Optional:    true,
			},
		},
	}
}

func (p *FirestoreProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	tflog.Info(ctx, "Configuring Firestore provider")

	var config FirestoreProviderModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Determine project
	project := os.Getenv("GOOGLE_PROJECT")
	if project == "" {
		project = os.Getenv("GOOGLE_CLOUD_PROJECT")
	}
	if !config.Project.IsNull() {
		project = config.Project.ValueString()
	}

	// Determine database
	database := "(default)"
	if !config.Database.IsNull() {
		database = config.Database.ValueString()
	}

	// Determine credentials
	credentials := os.Getenv("GOOGLE_CREDENTIALS")
	if credentials == "" {
		credentials = os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")
	}
	if !config.Credentials.IsNull() {
		credentials = config.Credentials.ValueString()
	}

	// Create HTTP client with authentication
	var httpClient *http.Client
	var err error

	scopes := []string{
		"https://www.googleapis.com/auth/datastore",
		"https://www.googleapis.com/auth/cloud-platform",
	}

	if credentials != "" {
		// Check if credentials is a file path or JSON content
		var credJSON []byte
		if _, statErr := os.Stat(credentials); statErr == nil {
			credJSON, err = os.ReadFile(credentials)
			if err != nil {
				resp.Diagnostics.AddError(
					"Unable to read credentials file",
					fmt.Sprintf("Error reading credentials file: %s", err),
				)
				return
			}
		} else {
			credJSON = []byte(credentials)
		}

		// Extract project from credentials if not set
		if project == "" {
			var credData map[string]interface{}
			if jsonErr := json.Unmarshal(credJSON, &credData); jsonErr == nil {
				if p, ok := credData["project_id"].(string); ok {
					project = p
				}
			}
		}

		creds, err := google.CredentialsFromJSON(context.Background(), credJSON, scopes...)
		if err != nil {
			resp.Diagnostics.AddError(
				"Unable to create credentials",
				fmt.Sprintf("Error creating credentials from JSON: %s", err),
			)
			return
		}

		httpClient, _, err = transport.NewHTTPClient(context.Background(), option.WithTokenSource(creds.TokenSource))
		if err != nil {
			resp.Diagnostics.AddError(
				"Unable to create HTTP client",
				fmt.Sprintf("Error creating HTTP client: %s", err),
			)
			return
		}
	} else {
		// Use Application Default Credentials
		creds, err := google.FindDefaultCredentials(context.Background(), scopes...)
		if err != nil {
			resp.Diagnostics.AddError(
				"Unable to find default credentials",
				fmt.Sprintf("Error finding default credentials: %s. Please set GOOGLE_APPLICATION_CREDENTIALS or configure credentials in the provider.", err),
			)
			return
		}

		if project == "" && creds.ProjectID != "" {
			project = creds.ProjectID
		}

		httpClient = oauth2.NewClient(context.Background(), creds.TokenSource)
	}

	if project == "" {
		resp.Diagnostics.AddError(
			"Missing project",
			"The provider could not determine the GCP project. Set the 'project' attribute or GOOGLE_PROJECT environment variable.",
		)
		return
	}

	client := &FirestoreClient{
		HTTPClient: httpClient,
		Project:    project,
		Database:   database,
	}

	tflog.Info(ctx, "Configured Firestore client", map[string]interface{}{
		"project":  project,
		"database": database,
	})

	resp.DataSourceData = client
	resp.ResourceData = client
}

func (p *FirestoreProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewDocumentResource,
	}
}

func (p *FirestoreProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewDocumentDataSource,
		NewDocumentsDataSource,
	}
}
