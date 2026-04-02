package backend.mapper;

import backend.dto.AllTechDTO;
import backend.dto.InternshipDTO;
import backend.entity.InternshipEntity;
import org.mapstruct.BeanMapping;
import org.mapstruct.Mapper;
import org.mapstruct.Mapping;
import org.mapstruct.ReportingPolicy;

@Mapper(componentModel = "spring", unmappedTargetPolicy = ReportingPolicy.ERROR)
public interface InternshipMapper {

    @Mapping(target = "companyName", source = "company.companyName")
    InternshipDTO toDTO(InternshipEntity internship);

}
