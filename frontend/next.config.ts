import type { NextConfig } from "next";

const nextConfig: NextConfig = {
  // Silence workspace root warning
  turbopack: {
    root: process.cwd(),
  },

  images: {
    remotePatterns: [
      {
        protocol: "https",
        hostname: "**.cdninstagram.com",
      },
      {
        protocol: "https",
        hostname: "scontent.cdninstagram.com",
      },
      {
        protocol: "https",
        hostname: "scontent*.cdninstagram.com",
      },
    ],
  },

  // Proxy API requests to the backend during development
  async rewrites() {
    return [
      {
        source: "/api/:path*",
        destination: `${(process.env.NEXT_PUBLIC_API_URL || "http://localhost:4000/api/v1").replace(/\/api\/v1\/?$/, "")}/api/v1/:path*`,
      },
    ];
  },
};

export default nextConfig;
