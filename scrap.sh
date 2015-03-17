#!/bin/bash

/home/marin/code/go/src/github.com/marinhero/columnScraper/getColumns > news.txt
mail -s "Columnas de Hoy" marin.alcaraz@gmail.com < news.txt
