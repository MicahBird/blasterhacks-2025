#!/usr/bin/env python3
import http.server
import socketserver
import os
import random
import mimetypes
import shutil
from urllib.parse import urlparse, unquote

# --- Configuration ---
PORT = 8080
ADS_BASE_DIR = './ads'
# Add more video extensions if needed
VIDEO_EXTENSIONS = {'.mp4', '.mov', '.avi', '.mkv', '.webm', '.mpeg', '.mpg'}
# ---

class AdRequestHandler(http.server.BaseHTTPRequestHandler):
    """
    Custom request handler to serve random videos from subdirectories.
    """
    def _send_response(self, code, message, content_type="text/html"):
        """Helper to send common responses."""
        self.send_response(code)
        self.send_header('Content-type', content_type)
        self.send_header('Content-Length', str(len(message)))
        # Add CORS headers to allow requests from any origin (optional, useful for testing)
        self.send_header('Access-Control-Allow-Origin', '*')
        self.end_headers()
        self.wfile.write(message.encode('utf-8'))

    def _list_available_endpoints(self):
        """Generates HTML listing available ad folder endpoints."""
        try:
            folders = [d for d in os.listdir(ADS_BASE_DIR)
                       if os.path.isdir(os.path.join(ADS_BASE_DIR, d))]
        except FileNotFoundError:
            return "<html><body><h1>Error</h1><p>Base ads directory '{}' not found.</p></body></html>".format(ADS_BASE_DIR), 500
        except Exception as e:
             return "<html><body><h1>Error</h1><p>Could not list directories: {}</p></body></html>".format(e), 500

        if not folders:
            return "<html><body><h1>No Ad Campaigns Found</h1><p>No subdirectories found in '{}'.</p></body></html>".format(ADS_BASE_DIR), 200

        list_items = "".join(f'<li><a href="/{folder}">{folder}</a></li>' for folder in sorted(folders))
        html = f"""
        <html>
          <head><title>Available Ad Endpoints</title></head>
          <body>
            <h1>Available Ad Endpoints</h1>
            <p>Serving videos from subdirectories inside '{ADS_BASE_DIR}':</p>
            <ul>{list_items}</ul>
          </body>
        </html>
        """
        return html, 200

    def do_GET(self):
        """Handles GET requests."""
        # --- Basic Security/Path Cleaning ---
        # Prevent directory traversal attacks
        parsed_path = urlparse(self.path)
        # Decode URL encoding (e.g., %20 -> space) and remove leading/trailing slashes
        requested_folder = unquote(parsed_path.path).strip('/')

        if ".." in requested_folder or requested_folder.startswith('/'):
             self._send_response(400, "<html><body><h1>Bad Request</h1><p>Invalid path.</p></body></html>")
             return

        # --- Handle Root Request: List Endpoints ---
        if not requested_folder:
            html_content, status_code = self._list_available_endpoints()
            self._send_response(status_code, html_content)
            return

        # --- Handle Ad Folder Request ---
        target_folder_path = os.path.abspath(os.path.join(ADS_BASE_DIR, requested_folder))

        # Double-check we are still within the intended base directory
        if not target_folder_path.startswith(os.path.abspath(ADS_BASE_DIR)):
            self._send_response(403, "<html><body><h1>Forbidden</h1><p>Access denied.</p></body></html>")
            print(f"Warning: Attempted access outside base directory: {self.path}")
            return

        if not os.path.isdir(target_folder_path):
            self._send_response(404, f"<html><body><h1>Not Found</h1><p>Folder '{requested_folder}' not found.</p></body></html>")
            return

        try:
            # Find video files in the target directory
            potential_files = os.listdir(target_folder_path)
            video_files = [
                f for f in potential_files
                if os.path.isfile(os.path.join(target_folder_path, f)) and
                   os.path.splitext(f)[1].lower() in VIDEO_EXTENSIONS
            ]

            if not video_files:
                self._send_response(404, f"<html><body><h1>Not Found</h1><p>No video files found in folder '{requested_folder}'.</p></body></html>")
                return

            # Select a random video
            chosen_video = random.choice(video_files)
            video_path = os.path.join(target_folder_path, chosen_video)

            # --- Serve the video file ---
            mime_type, _ = mimetypes.guess_type(video_path)
            if mime_type is None:
                mime_type = 'application/octet-stream' # Default if type can't be guessed

            try:
                file_size = os.path.getsize(video_path)
                with open(video_path, 'rb') as f:
                    self.send_response(200)
                    self.send_header('Content-type', mime_type)
                    self.send_header('Content-Length', str(file_size))
                    # Suggest browser display inline if possible
                    self.send_header('Content-Disposition', f'inline; filename="{chosen_video}"')
                    # Add CORS header (optional, but often useful)
                    self.send_header('Access-Control-Allow-Origin', '*')
                    self.end_headers()
                    # Efficiently copy file content to the response
                    shutil.copyfileobj(f, self.wfile)
                print(f"Served: {requested_folder}/{chosen_video} ({mime_type})")

            except FileNotFoundError:
                 # Should not happen if listing worked, but handle just in case
                 print(f"Error: File disappeared before serving: {video_path}")
                 self._send_response(404, "<html><body><h1>Not Found</h1><p>Video file not found (it may have been removed).</p></body></html>")
            except Exception as e:
                print(f"Error serving file {video_path}: {e}")
                # Avoid sending potentially sensitive error details to client
                self._send_response(500, "<html><body><h1>Internal Server Error</h1><p>Could not serve the video file.</p></body></html>")

        except Exception as e:
            print(f"Error processing request for {requested_folder}: {e}")
            self._send_response(500, "<html><body><h1>Internal Server Error</h1><p>An error occurred processing the request.</p></body></html>")


def run_server(port=PORT):
    # Ensure the base ads directory exists
    if not os.path.isdir(ADS_BASE_DIR):
        print(f"Error: Base ads directory '{ADS_BASE_DIR}' not found.")
        print("Please create it and add subfolders with videos.")
        return

    # Use ThreadingTCPServer to handle multiple requests concurrently
    socketserver.TCPServer.allow_reuse_address = True # Allow quick restarts
    httpd = socketserver.ThreadingTCPServer(("", port), AdRequestHandler)

    print(f"Serving HTTP on port {port}...")
    print(f"Serving random videos from subfolders in: {os.path.abspath(ADS_BASE_DIR)}")
    print(f"Access http://localhost:{port}/ to see available endpoints.")
    print(f"Access http://localhost:{port}/<folder_name> to get a random video.")
    print("Press Ctrl+C to stop.")

    try:
        httpd.serve_forever()
    except KeyboardInterrupt:
        print("\nServer stopping.")
        httpd.shutdown()
        httpd.server_close()

if __name__ == "__main__":
    run_server(port=PORT)
