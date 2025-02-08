package com.handlelapp.model

import org.jetbrains.exposed.sql.Table

object Products : Table("products") {
    val id = varchar("id", 255)
    val title = varchar("title", 255)
    val subtitle = varchar("subtitle", 255).nullable()
    val category = varchar("category", 255)
    val subcategory = varchar("subcategory", 255).nullable()
    val onsale = bool("onsale").nullable()
    val descripton = text("description").nullable()
    val weight = varchar("weight", 255).nullable()
    val origincountry = varchar("origincountry", 255).nullable()
    val ingredients = text("ingredients").nullable()
    val vendor = varchar("vendor", 255).nullable()
    val brand = varchar("brand", 255).nullable()
    val size = varchar("size", 255).nullable()
    val unit = varchar("unit", 255).nullable()
    val unittype = varchar("unittype", 255).nullable()
    val allergens = text("allergens").nullable()
    val mayContainTracesOf = text("maycontaintracesof").nullable()
    val nutritionalContent = text("nutritionalcontent").nullable()
    val prices = text("prices").nullable()
    val images = text("images").nullable()
    val notes = varchar("notes", 255).nullable()
}