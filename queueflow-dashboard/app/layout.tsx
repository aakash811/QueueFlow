import "./globals.css";
import Sidebar from "@/components/sidebar";

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="en">
      <body>
        <div className="flex">
          <Sidebar />

          <main className="flex-1 p-6 bg-zinc-950 text-white min-h-screen">
            {children}
          </main>
        </div>
      </body>
    </html>
  );
}
