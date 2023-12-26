import { Photo as TPhoto } from "@core/domain/domain.ts";

export default function Photo(opts: { photo: TPhoto; delay?: number }) {
  const { photo } = opts;
  return (
    <figure
      id={photo.id}
      className="item"
      // @ts-ignore
      style={{ "--delay": `${opts.delay || "0"}s` }}
    >
      <img src={photo.url} alt={photo.title} />

      <figcaption className="details">
        <h2>{photo.title}</h2>
        <div className="tags">
          {photo.tags.map((t) => (
            <span key={t}>{t}</span>
          ))}
        </div>
        <div className="data">
          <span>
            {photo.createdAt.substring(0, 10).split("-").reverse().join("/")}
          </span>
        </div>
      </figcaption>
    </figure>
  );
}
