import { Photo as TPhoto } from "@core/domain/domain.ts";

export default function Photo(opts: { photo: TPhoto; delay?: number }) {
  const { photo } = opts;
  return (
    // @ts-ignore
    <figure className="item" style={{ "--delay": `${opts.delay || "0"}s` }}>
      <img src={photo.url} alt={photo.title} />

      <figcaption className="details">
        <h2>{photo.title}</h2>
        <div className="tags">
          {photo.tags.map((t) => (
            <span key={t}>{t}</span>
          ))}
        </div>
      </figcaption>
    </figure>
  );
}
