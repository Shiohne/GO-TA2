import {AfterViewInit, Component, OnInit, ViewChild} from '@angular/core';
import {MatPaginator} from '@angular/material/paginator';
import {MatTableDataSource} from '@angular/material/table';
import { Dataset } from '../model/Dataset';
import { DatasetService } from '../services/dataset.service';

@Component({
  selector: 'app-data-set',
  templateUrl: './data-set.component.html',
  styleUrls: ['./data-set.component.css']
})
export class DataSetComponent implements AfterViewInit {
  displayedColumns: string[] = ['edad','tipo', 'actividad', 'insumo', 'metodo'];
  datos: any[]= [];
  dataSource = new MatTableDataSource<Dataset>();
  
  constructor( private datasetService: DatasetService){}

  @ViewChild(MatPaginator) paginator!: MatPaginator;

  ngAfterViewInit() {
    this.dataSource.paginator = this.paginator;
    this.getDataset()
  }
  getDataset() {
    this.datasetService.getDataset().subscribe((response: any) => {
    this.datos = response
    this.dataSource.data = this.datos
    })
  }
}