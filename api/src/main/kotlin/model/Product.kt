package com.handlelapp.model

import kotlinx.serialization.Serializable

@Serializable
data class Product(
    val id: String,
    val title: String,
    val subtitle: String?,
    val category: String,
    val subcategory: String?,
    val onsale: Boolean?,
    val description: String?,
    val weight: String?,
    val origincountry: String?,
    val ingredients: String?,
    val vendor: String?,
    val brand: String?,
    val size: String?,
    val unit: String?,
    val unittype: String?,
    val allergens: String?,
    val mayContainTracesOf: String?,
    val nutritionalContent: String?,
    val prices: String?,
    val images: String?,
    val notes: String?
)