import type { products } from '@prisma/client';

export type ExtendedProduct = products & {
	prices: Price[];
	nutritionalcontent: NutritionalContent;
	images: Image;
};

export type Price = {
	url: string;
	price: number;
	store: string;
	unitprice: number;
	originalprice: number;
};

type NutritionalContent = {
	fat: string;
	salt: string;
	energy: string;
	sodium: string;
	starch: string;
	sugars: string;
	protein: string;
	calories: string;
	dietaryfiber: string;
	saturatedfat: string;
	carbohydrates: string;
	monounsaturatedfat: string;
	polyunsaturatedfat: string;
};

type Image = {
	large: string;
	medium: string;
	small: string;
};
