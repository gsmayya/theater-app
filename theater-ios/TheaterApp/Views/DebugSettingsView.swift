import SwiftUI

struct DebugSettingsView: View {
    @StateObject private var apiService = APIService.shared
    @Environment(\.dismiss) private var dismiss
    
    var body: some View {
        NavigationView {
            VStack(spacing: 20) {
                // Header
                VStack(spacing: 8) {
                    Image(systemName: "gear")
                        .font(.system(size: 60))
                        .foregroundColor(.blue)
                    
                    Text("Debug Settings")
                        .font(.title)
                        .fontWeight(.bold)
                    
                    Text("Configure app behavior for testing")
                        .font(.subheadline)
                        .foregroundColor(.secondary)
                        .multilineTextAlignment(.center)
                }
                .padding(.top, 20)
                
                // Settings
                VStack(spacing: 16) {
                    // Data Source Toggle
                    VStack(alignment: .leading, spacing: 12) {
                        Text("Data Source")
                            .font(.headline)
                            .fontWeight(.semibold)
                        
                        Toggle("Use Mock Data", isOn: Binding(
                            get: { apiService.isUsingMockData },
                            set: { apiService.setUseMockData($0) }
                        ))
                        .toggleStyle(SwitchToggleStyle(tint: .blue))
                        
                        Text(apiService.isUsingMockData ? 
                             "Using sample data for testing" : 
                             "Using real backend API")
                            .font(.caption)
                            .foregroundColor(.secondary)
                    }
                    .padding()
                    .background(Color(.systemBackground))
                    .cornerRadius(12)
                    .shadow(color: .black.opacity(0.1), radius: 2, x: 0, y: 1)
                    
                    // Mock Data Info
                    if apiService.isUsingMockData {
                        VStack(alignment: .leading, spacing: 8) {
                            Text("Sample Data Includes:")
                                .font(.subheadline)
                                .fontWeight(.semibold)
                            
                            VStack(alignment: .leading, spacing: 4) {
                                Text("• 6 sample theater shows")
                                Text("• Romeo & Juliet, Lion King, Hamilton")
                                Text("• Phantom of the Opera, Les Misérables")
                                Text("• A Midsummer Night's Dream")
                                Text("• Multiple show times per show")
                                Text("• Sample bookings and users")
                            }
                            .font(.caption)
                            .foregroundColor(.secondary)
                        }
                        .padding()
                        .background(Color.blue.opacity(0.1))
                        .cornerRadius(12)
                    }
                    
                    // Backend Info
                    VStack(alignment: .leading, spacing: 8) {
                        Text("Backend Configuration")
                            .font(.subheadline)
                            .fontWeight(.semibold)
                        
                        VStack(alignment: .leading, spacing: 4) {
                            Text("URL: http://localhost:8080/api/v1")
                            Text("Status: \(apiService.isUsingMockData ? "Mock Data" : "Real API")")
                        }
                        .font(.caption)
                        .foregroundColor(.secondary)
                    }
                    .padding()
                    .background(Color(.systemGroupedBackground))
                    .cornerRadius(12)
                }
                
                Spacer()
            }
            .padding()
            .navigationTitle("Debug Settings")
            .navigationBarTitleDisplayMode(.inline)
            .toolbar {
                ToolbarItem(placement: .navigationBarTrailing) {
                    Button("Done") {
                        dismiss()
                    }
                }
            }
        }
    }
}

#Preview {
    DebugSettingsView()
}
