import { createFileRoute } from "@tanstack/react-router";
import { Github, Globe, Heart, Shield, Zap } from "lucide-react";
import { Badge } from "../components/ui/badge";
import { Button } from "../components/ui/button";
import {
  Card,
  CardContent,
  CardHeader,
  CardTitle,
} from "../components/ui/card";

export const Route = createFileRoute("/about")({
  component: About,
});

function About() {
  const features = [
    {
      icon: Zap,
      title: "Universal Proxy",
      desc: "One endpoint for all major AI providers (OpenAI, Anthropic, Gemini, etc.)",
    },
    {
      icon: Shield,
      title: "Auth Management",
      desc: "Easily manage and rotate API keys for local and remote access",
    },
    {
      icon: Heart,
      title: "Open Source",
      desc: "Built with transparency and community in mind",
    },
  ];

  return (
    <div className="p-6 max-w-4xl mx-auto animate-fade-in">
      <div className="mb-8 text-center">
        <div className="inline-block p-3 bg-white border-4 border-black shadow-neobrutal mb-4">
          <img
            src="/logo.png"
            alt="Quotio Logo"
            className="w-16 h-16 object-cover"
          />
        </div>
        <h2 className="text-3xl font-black mb-1">Quotio</h2>
        <p className="text-sm text-gray-500 font-bold uppercase tracking-widest">
          The Ultimate AI Proxy Manager v1.0.0
        </p>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-3 gap-6 mb-12">
        {features.map((f, i) => (
          <Card key={i} className="text-center">
            <CardContent className="pt-6">
              <div className="w-12 h-12 bg-black text-white rounded-full flex items-center justify-center mx-auto mb-4 border-4 border-black shadow-neobrutal-sm">
                <f.icon className="w-6 h-6" />
              </div>
              <h3 className="font-black text-base mb-2">{f.title}</h3>
              <p className="text-xs text-gray-600 font-medium leading-relaxed">
                {f.desc}
              </p>
            </CardContent>
          </Card>
        ))}
      </div>

      <Card className="mb-8">
        <CardHeader>
          <CardTitle>The Mission</CardTitle>
        </CardHeader>
        <CardContent className="space-y-4">
          <p className="text-sm text-gray-700 leading-relaxed font-medium">
            Quotio was created to solve the "API Key Proliferation" problem. As
            more AI agents and tools emerge, managing dozens of API keys across
            multiple providers becomes a security and organizational nightmare.
          </p>
          <p className="text-sm text-gray-700 leading-relaxed font-medium">
            By providing a local, secure proxy that mimics the OpenAI and
            Anthropic APIs, Quotio lets you rotate provider keys in one place
            while keeping your agents pointed at a single, authenticated local
            endpoint.
          </p>
        </CardContent>
      </Card>

      <div className="flex flex-wrap gap-4 justify-center">
        <a
          href="https://github.com/QuangThai/quotio-electron-go"
          target="_blank"
          rel="noreferrer"
        >
          <Button variant="secondary" className="flex items-center gap-2">
            <Github className="w-4 h-4" />
            GitHub Repository
          </Button>
        </a>
        <Button variant="secondary" className="flex items-center gap-2">
          <Globe className="w-4 h-4" />
          Documentation
        </Button>
      </div>

      <div className="mt-12 pt-6 border-t-2 border-black flex items-center justify-between text-[10px] font-bold text-gray-400 uppercase tracking-widest">
        <span>&copy; 2025 QUOTIO TEAM</span>
        <div className="flex gap-4">
          <Badge variant="secondary">STABLE</Badge>
          <Badge variant="secondary">LOCAL-ONLY</Badge>
        </div>
      </div>
    </div>
  );
}
