interface Props {
  title: string;
  url: string;
}

export default function GrafanaPanel({ title, url }: Props) {
  return (
    <div className="bg-zinc-900 border border-zinc-800 rounded-2xl overflow-hidden">
      <div className="p-4 border-b border-zinc-800">
        <h2 className="text-lg font-semibold text-white">{title}</h2>
      </div>

      <iframe
        src={url}
        width="100%"
        height="300"
        style={{ border: "none" }}
        className="bg-white"
      />
    </div>
  );
}
