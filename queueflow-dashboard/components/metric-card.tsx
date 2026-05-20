import { Card, CardContent } from "@/components/ui/card";

interface Props {
  title: string;
  value: string | number;
}

export default function MetricCard({ title, value }: Props) {
  return (
    <Card className="bg-zinc-900 border-zinc-800 text-white">
      <CardContent className="p-6">
        <p className="text-sm text-zinc-400">{title}</p>

        <h2 className="text-3xl font-bold mt-2">{value}</h2>
      </CardContent>
    </Card>
  );
}
