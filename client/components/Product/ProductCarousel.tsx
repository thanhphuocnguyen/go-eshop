import { GenericListResponse } from '@/lib/types';
import { CategoryProductModel } from '@/lib/types/product';
import Carousel from '../Common/Carousel';
import ProductCard from './ProductCard';

const responsive = {
  desktop: {
    breakpoint: { max: 3000, min: 1024 },
    items: 5,
    slidesToSlide: 3, // optional, default to 1.
  },
  tablet: {
    breakpoint: { max: 1024, min: 464 },
    items: 2,
    slidesToSlide: 2, // optional, default to 1.
  },
  mobile: {
    breakpoint: { max: 464, min: 0 },
    items: 1,
    slidesToSlide: 1, // optional, default to 1.
  },
};

export default async function ProductCarousel({
  categoryID,
}: {
  categoryID: number;
}) {
  const productResp: GenericListResponse<CategoryProductModel> = await fetch(
    process.env.NEXT_API_URL + `/category/${categoryID}/products`
  ).then((res) => res.json());
  return (
    <Carousel
      draggable={false}
      showDots={true}
      responsive={responsive}
      ssr // means to render carousel on server-side.
      infinite={true}
      autoPlay
      autoPlaySpeed={3000}
      customTransition='all .5'
      transitionDuration={800}
      containerClass='carousel-container'
      removeArrowOnDeviceType={['tablet', 'mobile']}
      deviceType={'desktop'}
      dotListClass='custom-dot-list-style'
      itemClass='carousel-item-padding-40-px'
    >
      {productResp.data.map((product) => (
        <ProductCard
          key={product.id}
          image={product.image_url}
          name={product.name}
          priceFrom={product.price_from}
          ID={product.id}
          priceTo={product.price_to}
          rating={4.5}
        />
      ))}
    </Carousel>
  );
}
