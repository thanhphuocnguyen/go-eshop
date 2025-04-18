export interface ColumnDefinition<T> {
  headerName: string;
  type: 'text' | 'number' | 'date';
  field: string;
  width?: number;
  render?: (data: T) => React.ReactNode;
}
