package com.handlelapp.repository

import com.handlelapp.model.Product
import com.handlelapp.model.Products
import org.jetbrains.exposed.sql.Query
import org.jetbrains.exposed.sql.selectAll
import org.jetbrains.exposed.sql.transactions.transaction

class ProductsRepository {
    fun getAllProducts(): List<Product> {
        return transaction {
            Products.selectAll().map {
                Product(
                    id = it[Products.id],
                    title = it[Products.title],
                    subtitle = it[Products.subtitle],
                    category = it[Products.category],
                    subcategory = it[Products.subcategory],
                    onsale = it[Products.onsale],
                    description = it[Products.descripton],
                    weight = it[Products.weight],
                    origincountry = it[Products.origincountry],
                    ingredients = it[Products.ingredients],
                    vendor = it[Products.vendor],
                    brand = it[Products.brand],
                    size = it[Products.size],
                    unit = it[Products.unit],
                    unittype = it[Products.unittype],
                    allergens = it[Products.allergens],
                    mayContainTracesOf = it[Products.mayContainTracesOf],
                    nutritionalContent = it[Products.nutritionalContent],
                    prices = it[Products.prices],
                    images = it[Products.images],
                    notes = it[Products.notes]
                )
            }
        }
    }
}