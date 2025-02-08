package com.handlelapp

import com.handlelapp.model.Product
import com.handlelapp.repository.ProductsRepository
import io.ktor.serialization.kotlinx.json.*
import io.ktor.server.application.*
import io.ktor.server.netty.*
import io.ktor.server.plugins.contentnegotiation.*
import io.ktor.server.response.*
import io.ktor.server.routing.*
import kotlinx.serialization.json.Json

fun main(args: Array<String>) {
    EngineMain.main(args)
}

fun Application.module() {
    DatabaseFactory.init()
    configureRouting()

    fun Application.configureSerialization() {
        install(ContentNegotiation) {
            json(Json { prettyPrint = true })
        }
    }
    configureSerialization()

    val products = ProductsRepository().getAllProducts()
    routing {
        get("/") {
            call.respond(products)
        }
    }
}