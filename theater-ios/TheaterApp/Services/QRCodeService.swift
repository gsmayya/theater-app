import Foundation
import CoreImage
import SwiftUI

// MARK: - QR Code Service
class QRCodeService {
    static let shared = QRCodeService()
    
    private init() {}
    
    // MARK: - Generate QR Code
    func generateQRCode(from string: String, size: CGSize = CGSize(width: 200, height: 200)) -> UIImage? {
        guard let data = string.data(using: .utf8) else { return nil }
        
        let context = CIContext()
        guard let filter = CIFilter(name: "CIQRCodeGenerator") else { return nil }
        
        filter.setValue(data, forKey: "inputMessage")
        filter.setValue("H", forKey: "inputCorrectionLevel")
        
        guard let outputImage = filter.outputImage else { return nil }
        
        let scaleX = size.width / outputImage.extent.size.width
        let scaleY = size.height / outputImage.extent.size.height
        let scaledImage = outputImage.transformed(by: CGAffineTransform(scaleX: scaleX, y: scaleY))
        
        guard let cgImage = context.createCGImage(scaledImage, from: scaledImage.extent) else { return nil }
        
        return UIImage(cgImage: cgImage)
    }
    
    // MARK: - Generate QR Code for Booking
    func generateBookingQRCode(booking: Booking) -> UIImage? {
        let qrData = BookingQRData(
            bookingId: booking.bookingId,
            showId: booking.showId,
            customerName: booking.customerName ?? "Guest",
            numberOfTickets: booking.numberOfTickets,
            bookingDate: booking.bookingDate,
            totalAmount: booking.totalAmount
        )
        
        let jsonString = qrData.toJSONString()
        return generateQRCode(from: jsonString)
    }
    
    // MARK: - Generate QR Code for Show
    func generateShowQRCode(show: Show) -> UIImage? {
        let qrData = ShowQRData(
            showId: show.id,
            title: show.title,
            location: show.location,
            price: show.price,
            showDate: show.showDate
        )
        
        let jsonString = qrData.toJSONString()
        return generateQRCode(from: jsonString)
    }
}

// MARK: - QR Code Data Models
struct BookingQRData: Codable {
    let bookingId: String
    let showId: String
    let customerName: String
    let numberOfTickets: Int
    let bookingDate: String
    let totalAmount: Int
    
    func toJSONString() -> String {
        guard let data = try? JSONEncoder().encode(self),
              let jsonString = String(data: data, encoding: .utf8) else {
            return ""
        }
        return jsonString
    }
}

struct ShowQRData: Codable {
    let showId: String
    let title: String
    let location: String
    let price: Int
    let showDate: String
    
    func toJSONString() -> String {
        guard let data = try? JSONEncoder().encode(self),
              let jsonString = String(data: data, encoding: .utf8) else {
            return ""
        }
        return jsonString
    }
}

// MARK: - QR Code View
struct QRCodeView: View {
    let qrImage: UIImage?
    let title: String
    let subtitle: String?
    
    var body: some View {
        VStack(spacing: 16) {
            if let qrImage = qrImage {
                Image(uiImage: qrImage)
                    .interpolation(.none)
                    .resizable()
                    .scaledToFit()
                    .frame(width: 200, height: 200)
                    .background(Color.white)
                    .cornerRadius(12)
                    .shadow(color: .black.opacity(0.1), radius: 4, x: 0, y: 2)
            } else {
                RoundedRectangle(cornerRadius: 12)
                    .fill(Color.gray.opacity(0.3))
                    .frame(width: 200, height: 200)
                    .overlay(
                        VStack {
                            Image(systemName: "qrcode")
                                .font(.system(size: 40))
                                .foregroundColor(.gray)
                            Text("QR Code Error")
                                .font(.caption)
                                .foregroundColor(.gray)
                        }
                    )
            }
            
            VStack(spacing: 4) {
                Text(title)
                    .font(.headline)
                    .fontWeight(.semibold)
                    .multilineTextAlignment(.center)
                
                if let subtitle = subtitle {
                    Text(subtitle)
                        .font(.subheadline)
                        .foregroundColor(.secondary)
                        .multilineTextAlignment(.center)
                }
            }
        }
        .padding()
    }
}

// MARK: - QR Code Scanner (Placeholder)
struct QRCodeScannerView: View {
    @Binding var scannedCode: String?
    @Environment(\.dismiss) private var dismiss
    
    var body: some View {
        VStack {
            Text("QR Code Scanner")
                .font(.title)
                .fontWeight(.bold)
                .padding()
            
            Text("Scanner functionality would be implemented here using AVFoundation")
                .font(.body)
                .foregroundColor(.secondary)
                .multilineTextAlignment(.center)
                .padding()
            
            Button("Dismiss") {
                dismiss()
            }
            .buttonStyle(.borderedProminent)
            
            Spacer()
        }
        .padding()
    }
}

#Preview {
    QRCodeView(
        qrImage: QRCodeService.shared.generateQRCode(from: "Sample QR Code Data"),
        title: "Sample QR Code",
        subtitle: "This is a preview"
    )
}
