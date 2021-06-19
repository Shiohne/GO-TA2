import {Component, OnInit} from '@angular/core';
import {FormBuilder, FormControl, FormGroup, Validators} from "@angular/forms";
import {MatTable, MatTableDataSource} from "@angular/material/table";
import {Usuaria} from "../model/Usuaria";
import {KnnService} from "../services/knn.service";
import {MatDialog} from "@angular/material/dialog";
import {DialogResponseComponent} from "../dialog-response/dialog-response.component";


@Component({
  selector: 'app-formulario',
  templateUrl: './formulario.component.html',
  styleUrls: ['./formulario.component.css']
})
export class FormularioComponent implements OnInit {

  dataSource = new MatTableDataSource<Usuaria>();
  datos: any[]= [];
  form!: FormGroup;
  prediccion: any;


  displayedColumns: string[] = [ 'edad', 'tipo', 'actividad', 'insumo','borrar'];
  constructor(private fb:FormBuilder, private knnService: KnnService,  private dialog: MatDialog ) {

  }

  ngOnInit(): void {
    this.datos = []
    this.form = new FormGroup({
      edad: new FormControl( '',[Validators.required, Validators.pattern(/^[1-9]\d{0,2}$/), Validators.max(80)]),
      tipo: new FormControl('',[ Validators.required, Validators.pattern(/^[0-9]\d{0,1}$/), Validators.max(1)]),
      actividad: new FormControl('',[Validators.required,Validators.pattern(/^[1-9]\d{0,10000}$/), Validators.max(2000)]),
      insumo: new FormControl('',[ Validators.required, Validators.pattern(/^[1-9]\d{0,100000}$/), Validators.max(10500)])
    })
  }

  addItem() {
    this.form.value.edad = parseFloat(this.form.value.edad)
    this.form.value.tipo = parseFloat(this.form.value.tipo)
    this.form.value.insumo = parseFloat(this.form.value.insumo)
    this.form.value.actividad = parseFloat(this.form.value.actividad)
    if(this.form.valid){
      this.datos.push(this.form.value)
    }
    this.dataSource.data = this.datos
    this.form.reset()
    console.log(this.datos)
  }
  knnPrediction() {
    this.knnService.postKnn(this.datos).subscribe((response: any) => {
    this.prediccion = response;
    this.dialog.open(DialogResponseComponent, {data: {respuesta: response}})
    })
  }
  Remove(element: any) {
    this.datos.forEach(((value, index) => {
      if (value == element)
      {
        this.datos.splice(index,1)
      }
    }))
    this.dataSource.data = this.datos
  }
  get edad() { return this.form.value.edad; }
  get tipo() { return this.form.value.tipo; }
  get actividad() { return this.form.value.actividad; }
  get insumo() { return this.form.value.insumo; }
}

