#!/bin/bash

./getColumns > news.txt
mail -s "Columnas de Hoy" marin.alcaraz@gmail.com < news.txt
