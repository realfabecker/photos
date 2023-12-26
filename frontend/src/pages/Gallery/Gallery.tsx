import { useEffect, useState } from "react";
import { State, useAppDispatch, useAppSelector } from "@store/store.ts";
import { getActionLoadPhotoList } from "@store/photos/creators/photo.ts";
import Photo from "@pages/Gallery/Photo.tsx";
import { Photo as TPhoto } from "@core/domain/domain";
import { ActionStatus } from "@core/domain/domain.ts";

const ErrorSection = (opts: { message: string }) => {
  return (
    <>
      <h1>Galeria</h1>
      <section className="gallery">
        <p>{opts.message}</p>
      </section>
    </>
  );
};

const ActionsSection = (opts: {
  photos: State<TPhoto[]>;
  onClickLoadMore: () => void;
}) => {
  return (
    <div className="actions">
      <button
        disabled={opts.photos.status == ActionStatus.LOADING}
        onClick={opts.onClickLoadMore}
      >
        {opts.photos.status == ActionStatus.LOADING
          ? "Carregando..."
          : "Carregar mais fotos"}
      </button>
    </div>
  );
};

export default function Gallery() {
  const dispatch = useAppDispatch();
  const [params, setParams] = useState({ page: 1, limit: 3 });
  const photos = useAppSelector((state) => state.photos["photos/list"]);

  useEffect(() => {
    dispatch(
      getActionLoadPhotoList({
        page: params.page,
        limit: params.limit,
        token: photos.page_token,
      })
    );
  }, [dispatch, params]);

  function handleClickLoadMore() {
    setParams((prevState) => ({
      ...prevState,
      page: prevState.page + 1,
    }));
  }

  if (photos.status === ActionStatus.ERROR) {
    return <ErrorSection message={photos?.error?.message || "Ops!"} />;
  }

  let k = 0;
  return (
    <>
      <h1>Galeria</h1>

      <section className="gallery">
        {(photos.data || []).map((p, i) => {
          k = i === 0 ? 0 : k > 0 && k % (params.limit - 1) === 0 ? 0 : k + 1;
          const delay = Math.floor((0.4 + k * 0.2) * 10) / 10;
          return <Photo key={p.id} photo={p} delay={delay} />;
        })}
      </section>

      {photos.has_more && (
        <ActionsSection photos={photos} onClickLoadMore={handleClickLoadMore} />
      )}
    </>
  );
}
